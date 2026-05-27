package transaction

import (
	"context"

	"github.com/tyzerrr/aws-log-practice/server/internal/adapter/db/sqlc"
)

type Tx interface {
	ReadOnlyTransaction(ctx context.Context, fn func(context.Context, sqlc.DBTX) error) error
	ReadWriteTransaction(ctx context.Context, fn func(context.Context, sqlc.DBTX) error) error
}
