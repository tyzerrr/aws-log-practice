package errors

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// status code mapping follows those between gRPC and HTTP.
const (
	CodeNotFoundErr int = 404
	CodeAlreadyExistsErr int = 409
	CodeInternalServerErr int = 500
)

// PostgreSQL SQLSTATE codes (https://www.postgresql.org/docs/current/errcodes-appendix.html).
const (
	sqlStateUniqueViolation = "23505"
)

type Error interface {
	Error() string
	Unwrap() error
	Code() int
	Retryable() bool
}

type NotFoundError struct {
	code int
	err error
}

func NewNotFoundError(err error) *NotFoundError {
	return &NotFoundError{
		code: CodeNotFoundErr,
		err: err,
	}
}

func (e *NotFoundError) Error() string {
	return e.err.Error()
}

func (e *NotFoundError) Unwrap() error {
	return e.err
}

func (e *NotFoundError) Code() int {
	return e.code
}

func (e *NotFoundError) Retryable() bool {
	return false
}

type AlreadyExistsError struct {
	code int
	err error
}

func NewAlreadyExistsError(err error) *AlreadyExistsError {
	return &AlreadyExistsError{
		code: CodeAlreadyExistsErr,
		err: err,
	}
}

func (e *AlreadyExistsError) Error() string {
	return e.err.Error()
}

func (e *AlreadyExistsError) Unwrap() error {
	return e.err
}

func (e *AlreadyExistsError) Code() int {
	return e.code
}

func (e *AlreadyExistsError) Retryable() bool {
	return false
}

type InternalServerError struct {
	code int
	err error
}

func NewInternalServerError(err error) *InternalServerError {
	return &InternalServerError{
		code: CodeInternalServerErr,
		err: err,
	}
}

func (e *InternalServerError) Error() string {
	return e.err.Error()
}

func (e *InternalServerError) Unwrap() error {
	return e.err
}

func (e *InternalServerError) Code() int {
	return e.code
}

func (e *InternalServerError) Retryable() bool {
	return true
}

// ErrorFromDB は pgx/PostgreSQL ドライバが返すエラーを domain error に変換する。
// sqlc 自体はステータスコードを持たず、ドライバのエラーをそのまま返すため、
// ここで pgx.ErrNoRows と pgconn.PgError の SQLSTATE を見て分類する。
func ErrorFromDB(err error) Error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return NewNotFoundError(err)
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case sqlStateUniqueViolation:
			return NewAlreadyExistsError(err)
		}
	}

	return NewInternalServerError(err)
}

