-- name: CreateOrder :one
INSERT INTO orders (
    transaction_id,
    product_id,
    quantity,
    unit_price_amount,
    total_price_amount,
    currency_code,
    status
) VALUES (
    sqlc.arg(transaction_id),
    sqlc.arg(product_id),
    sqlc.arg(quantity),
    sqlc.arg(unit_price_amount),
    sqlc.arg(total_price_amount),
    sqlc.arg(currency_code),
    sqlc.arg(status)
)
RETURNING *;

-- name: GetOrder :one
SELECT *
FROM orders
WHERE id = sqlc.arg(id);

-- name: ListOrdersByTransactionID :many
SELECT *
FROM orders
WHERE transaction_id = sqlc.arg(transaction_id)
ORDER BY ordered_at ASC;

-- name: ListOrdersByStatus :many
SELECT *
FROM orders
WHERE status = sqlc.arg(status)
ORDER BY ordered_at DESC
LIMIT sqlc.arg(limit_count)
OFFSET sqlc.arg(offset_count);

-- name: UpdateOrderStatus :one
UPDATE orders
SET
    status = sqlc.arg(status),
    updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: ListOrdersForProduct :many
SELECT *
FROM orders
WHERE product_id = sqlc.arg(product_id)
ORDER BY ordered_at DESC
LIMIT sqlc.arg(limit_count)
OFFSET sqlc.arg(offset_count);
