package main

import (
	"strings"
	"testing"
)

func TestDatabaseURLFromEnvPrefersDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://local-user:local-pass@localhost:5432/local-db?sslmode=disable")
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_NAME", "app_db")
	t.Setenv("DB_USER", "app-user")
	t.Setenv("DB_PASSWORD", "app-pass")

	got, err := databaseURLFromEnv()
	if err != nil {
		t.Fatalf("databaseURLFromEnv: %v", err)
	}

	want := "postgres://local-user:local-pass@localhost:5432/local-db?sslmode=disable"
	if got != want {
		t.Fatalf("databaseURLFromEnv() = %q, want %q", got, want)
	}
}

func TestDatabaseURLFromEnvBuildsURLFromDBEnv(t *testing.T) {
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "6543")
	t.Setenv("DB_NAME", "app_db")
	t.Setenv("DB_USER", "app-user")
	t.Setenv("DB_PASSWORD", "p@ss/word")
	t.Setenv("DB_SSLMODE", "verify-full")

	got, err := databaseURLFromEnv()
	if err != nil {
		t.Fatalf("databaseURLFromEnv: %v", err)
	}

	want := "postgres://app-user:p%40ss%2Fword@db.example.com:6543/app_db?sslmode=verify-full"
	if got != want {
		t.Fatalf("databaseURLFromEnv() = %q, want %q", got, want)
	}
}

func TestDatabaseURLFromEnvDefaultsPortAndSSLMode(t *testing.T) {
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_NAME", "app_db")
	t.Setenv("DB_USER", "app-user")
	t.Setenv("DB_PASSWORD", "app-pass")

	got, err := databaseURLFromEnv()
	if err != nil {
		t.Fatalf("databaseURLFromEnv: %v", err)
	}

	want := "postgres://app-user:app-pass@db.example.com:5432/app_db?sslmode=require"
	if got != want {
		t.Fatalf("databaseURLFromEnv() = %q, want %q", got, want)
	}
}

func TestDatabaseURLFromEnvRequiresDBEnv(t *testing.T) {
	t.Setenv("DB_NAME", "app_db")
	t.Setenv("DB_USER", "app-user")
	t.Setenv("DB_PASSWORD", "app-pass")

	_, err := databaseURLFromEnv()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "DB_HOST is required") {
		t.Fatalf("error = %q, want DB_HOST required", err.Error())
	}
}
