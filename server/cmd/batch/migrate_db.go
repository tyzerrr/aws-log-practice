package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultAtlasPath    = "/usr/local/bin/atlas"
	defaultMigrationDir = "/app/db/migrations"
	defaultShadowDBURL  = "postgres://shadow:shadow@127.0.0.1:5432/shadow?sslmode=disable"
	defaultDevDBURL     = "postgres://shadow:shadow@127.0.0.1:5432/shadow_dev?sslmode=disable"
)

type dbMigrationConfig struct {
	AWSRegion               string
	DBAdminCredentialID     string
	DBHost                  string
	DBPort                  string
	DBName                  string
	DBSSLMode               string
	Schema                  string
	AtlasPath               string
	MigrationDir            string
	MigrationBaseVersion    string
	ShadowDatabaseURL       string
	DevDatabaseURL          string
	ShadowMaintenanceDBName string
}

func CheckDBMigrationDrift(ctx context.Context) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return handleCheckDBMigrationDrift(ctx, logger)
}

func ApplyDBMigration(ctx context.Context) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return handleApplyDBMigration(ctx, logger)
}

func handleCheckDBMigrationDrift(ctx context.Context, logger *slog.Logger) error {
	cfg, err := loadDBMigrationConfig()
	if err != nil {
		return err
	}

	dbAdminCredential, err := getDBAdminCredential(ctx, cfg, logger)
	if err != nil {
		return err
	}
	targetURL := dbMigrationTargetURL(cfg, dbAdminCredential)

	if err := ensurePostgresDatabase(ctx, cfg.DevDatabaseURL, cfg.ShadowMaintenanceDBName); err != nil {
		return fmt.Errorf("failed to prepare atlas dev database: %w", err)
	}

	if cfg.MigrationBaseVersion != "" {
		if output, err := runAtlas(ctx, cfg.AtlasPath,
			"migrate",
			"apply",
			"--url", cfg.ShadowDatabaseURL,
			"--dir", migrationDirURL(cfg.MigrationDir),
			"--to-version", cfg.MigrationBaseVersion,
		); err != nil {
			return fmt.Errorf("failed to apply base migrations to shadow database: %w: %s", err, strings.TrimSpace(output))
		}
		logger.Info(
			"applied base migrations to shadow database",
			slog.String("base_version", cfg.MigrationBaseVersion),
		)
	} else {
		logger.Info("skipped base migration apply because no previous migration version exists")
	}

	diff, err := runAtlas(ctx, cfg.AtlasPath,
		"schema",
		"diff",
		"--from", cfg.ShadowDatabaseURL,
		"--to", targetURL,
		"--dev-url", cfg.DevDatabaseURL,
		"--schema", cfg.Schema,
		"--exclude", "atlas_schema_revisions",
		"--exclude", "*.atlas_schema_revisions",
		"--format", "{{ sql . }}",
	)
	if err != nil {
		return fmt.Errorf("failed to diff shadow database and remote database: %w: %s", err, strings.TrimSpace(diff))
	}
	if strings.TrimSpace(diff) != "" {
		logger.Error(
			"database schema drift detected",
			slog.String("db_name", cfg.DBName),
			slog.String("schema", cfg.Schema),
			slog.String("diff", strings.TrimSpace(diff)),
		)
		return fmt.Errorf("database schema drift detected")
	}

	logger.Info(
		"checked db migration drift",
		slog.String("db_name", cfg.DBName),
		slog.String("schema", cfg.Schema),
		slog.String("base_version", cfg.MigrationBaseVersion),
	)

	return nil
}

func handleApplyDBMigration(ctx context.Context, logger *slog.Logger) error {
	cfg, err := loadDBMigrationConfig()
	if err != nil {
		return err
	}

	dbAdminCredential, err := getDBAdminCredential(ctx, cfg, logger)
	if err != nil {
		return err
	}
	targetURL := dbMigrationTargetURL(cfg, dbAdminCredential)

	output, err := runAtlas(ctx, cfg.AtlasPath,
		"migrate",
		"apply",
		"--url", targetURL,
		"--dir", migrationDirURL(cfg.MigrationDir),
	)
	if err != nil {
		return fmt.Errorf("failed to apply db migrations: %w: %s", err, strings.TrimSpace(output))
	}

	logger.Info(
		"applied db migrations",
		slog.String("db_name", cfg.DBName),
		slog.String("schema", cfg.Schema),
		slog.String("atlas_output", strings.TrimSpace(output)),
	)

	return nil
}

func getDBAdminCredential(ctx context.Context, cfg dbMigrationConfig, logger *slog.Logger) (credential, error) {
	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.AWSRegion))
	if err != nil {
		return credential{}, fmt.Errorf("failed to load default aws config: %w", err)
	}

	smc := secretsmanager.NewFromConfig(awsConfig)
	client := &awsClient{
		smc:    smc,
		logger: logger,
	}

	dbAdminCredential, err := getCredential(ctx, client, cfg.DBAdminCredentialID)
	if err != nil {
		return credential{}, fmt.Errorf("failed to get db admin credential: %w", err)
	}

	return dbAdminCredential, nil
}

func loadDBMigrationConfig() (dbMigrationConfig, error) {
	cfg := dbMigrationConfig{
		AWSRegion:               strings.TrimSpace(os.Getenv("AWS_REGION")),
		DBAdminCredentialID:     strings.TrimSpace(os.Getenv("DB_ADMIN_CREDENTIAL_ID")),
		DBHost:                  firstEnv("DB_HOST", "DB_PRIMARY_HOST"),
		DBPort:                  envDefault("DB_PORT", defaultDBPort),
		DBName:                  strings.TrimSpace(os.Getenv("DB_NAME")),
		DBSSLMode:               envDefault("DB_SSLMODE", defaultDBSSLMode),
		Schema:                  envDefault("DB_SCHEMA", defaultSchema),
		AtlasPath:               envDefault("ATLAS_PATH", defaultAtlasPath),
		MigrationDir:            envDefault("MIGRATION_DIR", defaultMigrationDir),
		MigrationBaseVersion:    strings.TrimSpace(os.Getenv("MIGRATION_BASE_VERSION")),
		ShadowDatabaseURL:       envDefault("SHADOW_DATABASE_URL", defaultShadowDBURL),
		DevDatabaseURL:          envDefault("DEV_DATABASE_URL", defaultDevDBURL),
		ShadowMaintenanceDBName: envDefault("SHADOW_MAINTENANCE_DB_NAME", "postgres"),
	}

	switch {
	case cfg.AWSRegion == "":
		return dbMigrationConfig{}, fmt.Errorf("AWS_REGION is required")
	case cfg.DBAdminCredentialID == "":
		return dbMigrationConfig{}, fmt.Errorf("DB_ADMIN_CREDENTIAL_ID is required")
	case cfg.DBHost == "":
		return dbMigrationConfig{}, fmt.Errorf("DB_HOST is required")
	case cfg.DBName == "":
		return dbMigrationConfig{}, fmt.Errorf("DB_NAME is required")
	case cfg.Schema == "":
		return dbMigrationConfig{}, fmt.Errorf("DB_SCHEMA is required")
	case cfg.AtlasPath == "":
		return dbMigrationConfig{}, fmt.Errorf("ATLAS_PATH is required")
	case cfg.MigrationDir == "":
		return dbMigrationConfig{}, fmt.Errorf("MIGRATION_DIR is required")
	case cfg.ShadowDatabaseURL == "":
		return dbMigrationConfig{}, fmt.Errorf("SHADOW_DATABASE_URL is required")
	case cfg.DevDatabaseURL == "":
		return dbMigrationConfig{}, fmt.Errorf("DEV_DATABASE_URL is required")
	case cfg.ShadowMaintenanceDBName == "":
		return dbMigrationConfig{}, fmt.Errorf("SHADOW_MAINTENANCE_DB_NAME is required")
	}

	return cfg, nil
}

func dbMigrationTargetURL(cfg dbMigrationConfig, c credential) string {
	return postgresURL(createDBAppUserConfig{
		DBHost:    cfg.DBHost,
		DBPort:    cfg.DBPort,
		DBName:    cfg.DBName,
		DBSSLMode: cfg.DBSSLMode,
	}, c)
}

func migrationDirURL(dir string) string {
	if strings.HasPrefix(dir, "file://") {
		return dir
	}
	if filepath.IsAbs(dir) {
		return "file://" + dir
	}
	return "file://" + filepath.ToSlash(dir)
}

func runAtlas(ctx context.Context, atlasPath string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, atlasPath, args...)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		return output.String(), err
	}
	return output.String(), nil
}

func ensurePostgresDatabase(ctx context.Context, databaseURL string, maintenanceDBName string) error {
	targetURL, err := url.Parse(databaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse database url: %w", err)
	}
	targetDBName := strings.TrimPrefix(targetURL.Path, "/")
	if targetDBName == "" {
		return fmt.Errorf("database url path must include database name")
	}

	maintenanceURL := *targetURL
	maintenanceURL.Path = "/" + maintenanceDBName

	pool, err := pgxpool.New(ctx, maintenanceURL.String())
	if err != nil {
		return fmt.Errorf("failed to connect to maintenance database: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping maintenance database: %w", err)
	}

	var exists bool
	if err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", targetDBName).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}
	if exists {
		return nil
	}

	_, err = pool.Exec(ctx, "CREATE DATABASE "+pgxIdentifier(targetDBName))
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}

func pgxIdentifier(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}
