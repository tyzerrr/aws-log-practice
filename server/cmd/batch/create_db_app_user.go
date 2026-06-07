package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/jackc/pgx/v5"
)

const (
	versionStage     = "AWSCURRENT"
	defaultDBPort    = "5432"
	defaultDBSSLMode = "require"
	defaultSchema    = "public"
)

//go:embed sql/create_db_app_user.sql.tmpl
var createDBAppUserSQLTemplate string

type Client interface {
	GetSecretValue(ctx context.Context, input *secretsmanager.GetSecretValueInput) ([]byte, error)
}

type credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type createDBAppUserConfig struct {
	AWSRegion           string
	DBAdminCredentialID string
	DBAppCredentialID   string
	DBHost              string
	DBPort              string
	DBName              string
	DBSSLMode           string
	Schema              string
}

type createDBAppUserSQLData struct {
	UsernameIdent   string
	UsernameLiteral string
	PasswordLiteral string
	DatabaseIdent   string
	SchemaIdent     string
}

func CreateDBAppUser(ctx context.Context) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return handleCreateDBAppUser(ctx, logger)
}

func handleCreateDBAppUser(ctx context.Context, logger *slog.Logger) error {
	cfg, err := loadCreateDBAppUserConfig()
	if err != nil {
		return err
	}

	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.AWSRegion))
	if err != nil {
		return fmt.Errorf("failed to load default aws config: %w", err)
	}

	smc := secretsmanager.NewFromConfig(awsConfig)
	client := &awsClient{
		smc:    smc,
		logger: logger,
	}

	dbAdminCredential, err := getCredential(ctx, client, cfg.DBAdminCredentialID)
	if err != nil {
		return fmt.Errorf("failed to get db admin credential: %w", err)
	}

	dbAppCredential, err := getCredential(ctx, client, cfg.DBAppCredentialID)
	if err != nil {
		return fmt.Errorf("failed to get db app credential: %w", err)
	}

	database, err := newDB(ctx, postgresURL(cfg, dbAdminCredential))
	if err != nil {
		return err
	}
	defer database.Close()

	sql, err := renderCreateDBAppUserSQL(cfg, dbAppCredential)
	if err != nil {
		return err
	}

	if err := database.ReadWriteTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		return execSQL(ctx, tx, sql)
	}); err != nil {
		return err
	}

	logger.Info(
		"created or updated db app user",
		slog.String("db_name", cfg.DBName),
		slog.String("schema", cfg.Schema),
		slog.String("username", dbAppCredential.Username),
	)

	return nil
}

func execSQL(ctx context.Context, tx pgx.Tx, sql string) error {
	if _, err := tx.Conn().PgConn().Exec(ctx, sql).ReadAll(); err != nil {
		return fmt.Errorf("failed to execute sql: %w", err)
	}
	return nil
}

func getCredential(ctx context.Context, client Client, secretID string) (credential, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretID),
		VersionStage: aws.String(versionStage), // VersionStage defaults to AWSCURRENT if unspecified.
	}

	raw, err := client.GetSecretValue(ctx, input)
	if err != nil {
		return credential{}, fmt.Errorf("failed to get secret value from SecretsManager: %w", err)
	}

	var c credential
	if err := json.Unmarshal(raw, &c); err != nil {
		return credential{}, fmt.Errorf("failed to unmarshal credential: %w", err)
	}
	if c.Username == "" {
		return credential{}, fmt.Errorf("credential username is empty")
	}
	if c.Password == "" {
		return credential{}, fmt.Errorf("credential password is empty")
	}

	return c, nil
}

func loadCreateDBAppUserConfig() (createDBAppUserConfig, error) {
	cfg := createDBAppUserConfig{
		AWSRegion:           strings.TrimSpace(os.Getenv("AWS_REGION")),
		DBAdminCredentialID: strings.TrimSpace(os.Getenv("DB_ADMIN_CREDENTIAL_ID")),
		DBAppCredentialID:   strings.TrimSpace(os.Getenv("DB_APP_CREDENTIAL_ID")),
		DBHost:              firstEnv("DB_HOST", "DB_PRIMARY_HOST"),
		DBPort:              envDefault("DB_PORT", defaultDBPort),
		DBName:              strings.TrimSpace(os.Getenv("DB_NAME")),
		DBSSLMode:           envDefault("DB_SSLMODE", defaultDBSSLMode),
		Schema:              envDefault("DB_SCHEMA", defaultSchema),
	}

	switch {
	case cfg.AWSRegion == "":
		return createDBAppUserConfig{}, fmt.Errorf("AWS_REGION is required")
	case cfg.DBAdminCredentialID == "":
		return createDBAppUserConfig{}, fmt.Errorf("DB_ADMIN_CREDENTIAL_ID is required")
	case cfg.DBAppCredentialID == "":
		return createDBAppUserConfig{}, fmt.Errorf("DB_APP_CREDENTIAL_ID is required")
	case cfg.DBHost == "":
		return createDBAppUserConfig{}, fmt.Errorf("DB_HOST is required")
	case cfg.DBName == "":
		return createDBAppUserConfig{}, fmt.Errorf("DB_NAME is required")
	case cfg.Schema == "":
		return createDBAppUserConfig{}, fmt.Errorf("DB_SCHEMA is required")
	}

	return cfg, nil
}

func renderCreateDBAppUserSQL(cfg createDBAppUserConfig, appCredential credential) (string, error) {
	usernameLiteral, err := quoteSQLLiteral(appCredential.Username)
	if err != nil {
		return "", fmt.Errorf("failed to quote username: %w", err)
	}
	passwordLiteral, err := quoteSQLLiteral(appCredential.Password)
	if err != nil {
		return "", fmt.Errorf("failed to quote password: %w", err)
	}

	data := createDBAppUserSQLData{
		UsernameIdent:   pgx.Identifier{appCredential.Username}.Sanitize(),
		UsernameLiteral: usernameLiteral,
		PasswordLiteral: passwordLiteral,
		DatabaseIdent:   pgx.Identifier{cfg.DBName}.Sanitize(),
		SchemaIdent:     pgx.Identifier{cfg.Schema}.Sanitize(),
	}

	tmpl, err := template.New("create_db_app_user.sql.tmpl").Parse(createDBAppUserSQLTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse create db app user sql template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render create db app user sql template: %w", err)
	}

	return buf.String(), nil
}

func postgresURL(cfg createDBAppUserConfig, c credential) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.Username, c.Password),
		Host:   net.JoinHostPort(cfg.DBHost, cfg.DBPort),
		Path:   "/" + cfg.DBName,
	}
	values := url.Values{}
	values.Set("sslmode", cfg.DBSSLMode)
	u.RawQuery = values.Encode()

	return u.String()
}

func quoteSQLLiteral(value string) (string, error) {
	if strings.ContainsRune(value, 0) {
		return "", fmt.Errorf("value contains null byte")
	}
	return "'" + strings.ReplaceAll(value, "'", "''") + "'", nil
}

func envDefault(key string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}
