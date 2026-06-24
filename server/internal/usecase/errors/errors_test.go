package errors

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestErrorFromDB(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      error
		wantNil  bool
		wantCode int
	}{
		{
			name:    "nil returns nil",
			err:     nil,
			wantNil: true,
		},
		{
			name:     "no rows maps to not found",
			err:      pgx.ErrNoRows,
			wantCode: CodeNotFoundErr,
		},
		{
			name:     "unique violation maps to already exists",
			err:      &pgconn.PgError{Code: sqlStateUniqueViolation},
			wantCode: CodeAlreadyExistsErr,
		},
		{
			name:     "unknown pg error maps to internal server error",
			err:      &pgconn.PgError{Code: "08006"},
			wantCode: CodeInternalServerErr,
		},
		{
			name:     "non pg error maps to internal server error",
			err:      errors.New("boom"),
			wantCode: CodeInternalServerErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ErrorFromDB(tt.err)

			if tt.wantNil {
				if got != nil {
					t.Fatalf("ErrorFromDB(%v) = %v, want nil", tt.err, got)
				}
				return
			}

			if got == nil {
				t.Fatalf("ErrorFromDB(%v) = nil, want code %d", tt.err, tt.wantCode)
			}
			if got.Code() != tt.wantCode {
				t.Fatalf("ErrorFromDB(%v).Code() = %d, want %d", tt.err, got.Code(), tt.wantCode)
			}
			if !errors.Is(got, tt.err) {
				t.Fatalf("ErrorFromDB(%v) should wrap the original error", tt.err)
			}
		})
	}
}
