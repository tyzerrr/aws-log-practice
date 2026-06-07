package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/jackc/pgx/v5"
)

//go:embed sql/check_db_app_user.sql.tmpl
var checkDBAppUserSQLTemplate string

type checkDBAppUserSQLData struct {
	DatabaseLiteral string
	SchemaLiteral   string
}

type privilegeCheck struct {
	Username     string
	DatabaseName string
	TargetType   string
	TargetName   string
	Privilege    string
	OK           bool
}

func CheckDBAppUser(ctx context.Context) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return handleCheckDBAppUser(ctx, logger)
}

func handleCheckDBAppUser(ctx context.Context, logger *slog.Logger) error {
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

	dbAppCredential, err := getCredential(ctx, client, cfg.DBAppCredentialID)
	if err != nil {
		return fmt.Errorf("failed to get db app credential: %w", err)
	}

	database, err := newDB(ctx, postgresURL(cfg, dbAppCredential))
	if err != nil {
		return err
	}
	defer database.Close()

	sql, err := renderCheckDBAppUserSQL(cfg)
	if err != nil {
		return err
	}

	var checks []privilegeCheck
	if err := database.ReadOnlyTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		checks, err = queryPrivilegeChecks(ctx, tx, sql)
		return err
	}); err != nil {
		return err
	}

	failures := failedPrivilegeChecks(checks)
	if len(failures) > 0 {
		for _, failure := range failures {
			logger.Error(
				"db app user privilege check failed",
				slog.String("username", failure.Username),
				slog.String("db_name", failure.DatabaseName),
				slog.String("target_type", failure.TargetType),
				slog.String("target_name", failure.TargetName),
				slog.String("privilege", failure.Privilege),
			)
		}
		return fmt.Errorf("db app user privilege check failed: %d failed checks", len(failures))
	}

	logger.Info(
		"checked db app user privileges",
		slog.String("db_name", cfg.DBName),
		slog.String("schema", cfg.Schema),
		slog.String("username", dbAppCredential.Username),
		slog.Int("checks_count", len(checks)),
	)

	return nil
}

func renderCheckDBAppUserSQL(cfg createDBAppUserConfig) (string, error) {
	databaseLiteral, err := quoteSQLLiteral(cfg.DBName)
	if err != nil {
		return "", fmt.Errorf("failed to quote database name: %w", err)
	}
	schemaLiteral, err := quoteSQLLiteral(cfg.Schema)
	if err != nil {
		return "", fmt.Errorf("failed to quote schema name: %w", err)
	}

	data := checkDBAppUserSQLData{
		DatabaseLiteral: databaseLiteral,
		SchemaLiteral:   schemaLiteral,
	}

	tmpl, err := template.New("check_db_app_user.sql.tmpl").Parse(checkDBAppUserSQLTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse check db app user sql template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render check db app user sql template: %w", err)
	}

	return buf.String(), nil
}

func queryPrivilegeChecks(ctx context.Context, tx pgx.Tx, sql string) ([]privilegeCheck, error) {
	rows, err := tx.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to query db app user privilege checks: %w", err)
	}
	defer rows.Close()

	var checks []privilegeCheck
	for rows.Next() {
		var check privilegeCheck
		if err := rows.Scan(
			&check.Username,
			&check.DatabaseName,
			&check.TargetType,
			&check.TargetName,
			&check.Privilege,
			&check.OK,
		); err != nil {
			return nil, fmt.Errorf("failed to scan db app user privilege check: %w", err)
		}
		checks = append(checks, check)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate db app user privilege checks: %w", err)
	}

	return checks, nil
}

func failedPrivilegeChecks(checks []privilegeCheck) []privilegeCheck {
	failures := make([]privilegeCheck, 0)
	for _, check := range checks {
		if !check.OK {
			failures = append(failures, check)
		}
	}
	return failures
}
