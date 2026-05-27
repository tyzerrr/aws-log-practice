-- name: CreateTransaction :one
INSERT INTO transactions (
    total_price_amount,
    currency_code,
    status
) VALUES (
    sqlc.arg(total_price_amount),
    sqlc.arg(currency_code),
    sqlc.arg(status)
)
RETURNING *;

-- name: GetTransaction :one
SELECT *
FROM transactions
WHERE id = sqlc.arg(id);

-- name: ListTransactionsByStatus :many
SELECT *
FROM transactions
WHERE status = sqlc.arg(status)
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count)
OFFSET sqlc.arg(offset_count);

-- name: UpdateTransactionStatus :one
UPDATE transactions
SET
    status = sqlc.arg(status),
    updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateTransactionTotalPriceAmount :one
UPDATE transactions
SET
    total_price_amount = sqlc.arg(total_price_amount),
    updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;
