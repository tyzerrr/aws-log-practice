-- name: CreateStock :one
INSERT INTO stocks (
    product_id,
    quantity
) VALUES (
    sqlc.arg(product_id),
    sqlc.arg(quantity)
)
RETURNING *;

-- name: GetStockByProductID :one
SELECT *
FROM stocks
WHERE product_id = sqlc.arg(product_id);

-- name: UpdateStockQuantity :one
UPDATE stocks
SET
    quantity = sqlc.arg(quantity),
    updated_at = now()
WHERE product_id = sqlc.arg(product_id)
RETURNING *;

-- name: ReserveStock :one
UPDATE stocks
SET
    reserved_quantity = reserved_quantity + sqlc.arg(quantity),
    updated_at = now()
WHERE product_id = sqlc.arg(product_id)
  AND quantity - reserved_quantity >= sqlc.arg(quantity)
RETURNING *;

-- name: ReleaseReservedStock :one
UPDATE stocks
SET
    reserved_quantity = reserved_quantity - sqlc.arg(quantity),
    updated_at = now()
WHERE product_id = sqlc.arg(product_id)
  AND reserved_quantity >= sqlc.arg(quantity)
RETURNING *;

-- name: CommitReservedStock :one
UPDATE stocks
SET
    quantity = quantity - sqlc.arg(quantity),
    reserved_quantity = reserved_quantity - sqlc.arg(quantity),
    updated_at = now()
WHERE product_id = sqlc.arg(product_id)
  AND reserved_quantity >= sqlc.arg(quantity)
RETURNING *;

-- name: ListLowStocks :many
SELECT *
FROM stocks
WHERE quantity - reserved_quantity <= sqlc.arg(threshold)
ORDER BY quantity - reserved_quantity ASC, updated_at DESC;
