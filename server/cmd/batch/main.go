package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
	ExitOK int = iota
	ExitErr
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := realMain(ctx); err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		os.Exit(ExitErr)
	}
	os.Exit(ExitOK)
}

func realMain(ctx context.Context) error {
	if len(os.Args) != 2 {
		return fmt.Errorf("usage: go run ${ROOT of aws-log-practice}/server/cmd/batch/ [create_db_app_user|check_db_app_user|check_db_migration_drift|apply_db_migration]")
	}
	var err error

	switch os.Args[1] {
	case "create_db_app_user":
		err = CreateDBAppUser(ctx)
	case "check_db_app_user":
		err = CheckDBAppUser(ctx)
	case "check_db_migration_drift":
		err = CheckDBMigrationDrift(ctx)
	case "apply_db_migration":
		err = ApplyDBMigration(ctx)
	default:
		err = fmt.Errorf("invalid command: %s", os.Args[1])
	}

	return err
}
