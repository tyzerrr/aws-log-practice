package main

import (
	"strings"
	"testing"
)

func TestRenderCreateDBAppUserSQL(t *testing.T) {
	t.Parallel()

	sql, err := renderCreateDBAppUserSQL(
		createDBAppUserConfig{
			DBName: "app-db",
			Schema: "public",
		},
		credential{
			Username: "app-user",
			Password: "p'ass",
		},
	)
	if err != nil {
		t.Fatalf("render create db app user sql: %v", err)
	}

	for _, want := range []string{
		`CREATE ROLE "app-user" WITH LOGIN PASSWORD 'p''ass';`,
		`ALTER ROLE "app-user" WITH LOGIN PASSWORD 'p''ass';`,
		`GRANT CONNECT ON DATABASE "app-db" TO "app-user";`,
		`GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA "public" TO "app-user";`,
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("rendered sql does not contain %q:\n%s", want, sql)
		}
	}
}

func TestRenderCheckDBAppUserSQL(t *testing.T) {
	t.Parallel()

	sql, err := renderCheckDBAppUserSQL(createDBAppUserConfig{
		DBName: "app-db",
		Schema: "public",
	})
	if err != nil {
		t.Fatalf("render check db app user sql: %v", err)
	}

	for _, want := range []string{
		`has_database_privilege(current_user, 'app-db', 'CONNECT') AS ok`,
		`has_schema_privilege(current_user, 'public', 'USAGE') AS ok`,
		`has_table_privilege(current_user, format('%I.%I', schemaname, tablename), privilege) AS ok`,
		`has_sequence_privilege(current_user, format('%I.%I', schemaname, sequencename), privilege) AS ok`,
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("rendered sql does not contain %q:\n%s", want, sql)
		}
	}
}

func TestQuoteSQLLiteralRejectsNullByte(t *testing.T) {
	t.Parallel()

	if _, err := quoteSQLLiteral("a\x00b"); err == nil {
		t.Fatal("expected error")
	}
}
