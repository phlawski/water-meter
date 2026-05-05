-- name: ListReadings :many
SELECT id, read_at, value_m3, price_per_m3, notes
FROM readings
ORDER BY read_at DESC;

-- name: InsertReading :one
INSERT INTO readings (read_at, value_m3, price_per_m3, notes)
VALUES (?, ?, ?, ?)
RETURNING id, read_at, value_m3, price_per_m3, notes;

-- name: DeleteReading :exec
DELETE FROM readings WHERE id = ?;

-- name: GetReading :one
SELECT id, read_at, value_m3, price_per_m3, notes
FROM readings
WHERE id = ?;

-- name: GetConfigValue :one
SELECT value FROM config WHERE key = ?;

-- name: SetConfigValue :exec
INSERT INTO config (key, value) VALUES (?, ?)
ON CONFLICT(key) DO UPDATE SET value = excluded.value;
