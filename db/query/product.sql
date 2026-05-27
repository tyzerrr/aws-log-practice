-- name: CreateProduct :one
INSERT INTO products (
    name,
    description,
    price_amount,
    currency_code,
    is_active
) VALUES (
    sqlc.arg(name),
    sqlc.arg(description),
    sqlc.arg(price_amount),
    sqlc.arg(currency_code),
    sqlc.arg(is_active)
)
RETURNING *;

-- name: GetProduct :one
SELECT *
FROM products
WHERE id = sqlc.arg(id);

-- name: ListActiveProducts :many
SELECT *
FROM products
WHERE is_active = TRUE
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count)
OFFSET sqlc.arg(offset_count);

-- name: UpdateProduct :one
UPDATE products
SET
    name = sqlc.arg(name),
    description = sqlc.arg(description),
    price_amount = sqlc.arg(price_amount),
    currency_code = sqlc.arg(currency_code),
    is_active = sqlc.arg(is_active),
    updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeactivateProduct :one
UPDATE products
SET
    is_active = FALSE,
    updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;
