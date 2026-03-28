-- ============================================================================
-- Users Queries
-- ============================================================================

-- name: CreateAnonymousUser :one
INSERT INTO users (id, low_stock_threshold, notification_enabled)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateUser :one
UPDATE users
SET
    name = COALESCE(sqlc.narg('name'), name),
    low_stock_threshold = COALESCE(sqlc.narg('low_stock_threshold'), low_stock_threshold),
    notification_enabled = COALESCE(sqlc.narg('notification_enabled'), notification_enabled),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;

-- ============================================================================
-- Coffee Beans Queries (without purchase fields)
-- ============================================================================

-- name: CreateCoffeeBean :one
INSERT INTO coffee_beans (
    id, user_id, name, origin, roast_level, current_stock, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetCoffeeBeanByID :one
SELECT * FROM coffee_beans
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListCoffeeBeansByUserID :many
SELECT * FROM coffee_beans
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: ListCoffeeBeansWithLatestPurchase :many
SELECT
    cb.id,
    cb.user_id,
    cb.name,
    cb.origin,
    cb.roast_level,
    cb.current_stock,
    cb.notes,
    cb.created_at,
    cb.updated_at,
    cb.deleted_at,
    ph.purchase_date as latest_purchase_date,
    ph.purchase_price as latest_purchase_price,
    ph.purchase_store as latest_purchase_store
FROM coffee_beans cb
LEFT JOIN LATERAL (
    SELECT purchase_date, purchase_price, purchase_store
    FROM purchase_history
    WHERE coffee_bean_id = cb.id
    ORDER BY purchase_date DESC, created_at DESC
    LIMIT 1
) ph ON true
WHERE cb.user_id = $1 AND cb.deleted_at IS NULL
ORDER BY cb.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateCoffeeBean :one
UPDATE coffee_beans
SET
    name = COALESCE(sqlc.narg('name'), name),
    origin = COALESCE(sqlc.narg('origin'), origin),
    roast_level = COALESCE(sqlc.narg('roast_level'), roast_level),
    current_stock = COALESCE(sqlc.narg('current_stock'), current_stock),
    notes = COALESCE(sqlc.narg('notes'), notes),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateCoffeeBeanStock :one
UPDATE coffee_beans
SET
    current_stock = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteCoffeeBean :exec
UPDATE coffee_beans
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: CountCoffeeBeansByUserID :one
SELECT COUNT(*) FROM coffee_beans
WHERE user_id = $1 AND deleted_at IS NULL;

-- ============================================================================
-- Purchase History Queries (NEW)
-- ============================================================================

-- name: CreatePurchaseHistory :one
INSERT INTO purchase_history (
    coffee_bean_id, user_id, purchase_date, purchase_price,
    purchase_store, quantity, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetPurchaseHistoryByID :one
SELECT * FROM purchase_history
WHERE id = $1;

-- name: GetPurchaseHistoriesByCoffeeBeanID :many
SELECT * FROM purchase_history
WHERE coffee_bean_id = $1
ORDER BY purchase_date DESC, created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetLatestPurchaseHistoryByCoffeeBeanID :one
SELECT * FROM purchase_history
WHERE coffee_bean_id = $1
ORDER BY purchase_date DESC, created_at DESC
LIMIT 1;

-- name: ListPurchaseHistoriesByUserID :many
SELECT ph.*, cb.name as coffee_bean_name
FROM purchase_history ph
JOIN coffee_beans cb ON ph.coffee_bean_id = cb.id
WHERE ph.user_id = $1
ORDER BY ph.purchase_date DESC, ph.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePurchaseHistory :one
UPDATE purchase_history
SET
    purchase_date = COALESCE(sqlc.narg('purchase_date'), purchase_date),
    purchase_price = COALESCE(sqlc.narg('purchase_price'), purchase_price),
    purchase_store = COALESCE(sqlc.narg('purchase_store'), purchase_store),
    quantity = COALESCE(sqlc.narg('quantity'), quantity),
    notes = COALESCE(sqlc.narg('notes'), notes)
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeletePurchaseHistory :exec
DELETE FROM purchase_history
WHERE id = $1 AND user_id = $2;

-- name: CountPurchaseHistoriesByCoffeeBeanID :one
SELECT COUNT(*) FROM purchase_history
WHERE coffee_bean_id = $1;

-- ============================================================================
-- Usage History Queries
-- ============================================================================

-- name: CreateUsageHistory :one
INSERT INTO usage_history (
    coffee_bean_id, user_id, usage_date, quantity, usage_type, notes
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUsageHistoryByID :one
SELECT * FROM usage_history
WHERE id = $1;

-- name: ListUsageHistoriesByUserID :many
SELECT uh.*, cb.name as coffee_bean_name
FROM usage_history uh
JOIN coffee_beans cb ON uh.coffee_bean_id = cb.id
WHERE uh.user_id = $1 AND cb.deleted_at IS NULL
ORDER BY uh.usage_date DESC, uh.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListUsageHistoriesByCoffeeBeanID :many
SELECT * FROM usage_history
WHERE coffee_bean_id = $1
ORDER BY usage_date DESC, created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetRecentUsageHistoryForConsumptionRate :many
SELECT usage_date, quantity
FROM usage_history
WHERE coffee_bean_id = $1 AND usage_date >= $2
ORDER BY usage_date DESC;

-- name: DeleteUsageHistory :exec
DELETE FROM usage_history
WHERE id = $1 AND user_id = $2;

-- name: CountUsageHistoriesByUserID :one
SELECT COUNT(*) FROM usage_history
WHERE user_id = $1;

-- name: CountUsageHistoriesByCoffeeBeanID :one
SELECT COUNT(*) FROM usage_history
WHERE coffee_bean_id = $1;
