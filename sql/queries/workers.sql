-- name: GetWorker :one
SELECT * FROM workers
WHERE id = $1;

-- name: GetWorkers :many
SELECT * FROM workers;

-- name: CreateWorker :one
INSERT INTO workers (
    host, ip_addr, connected_at, last_seen_at
) VALUES (
    $1, $2, NOW(), NOW()
)
RETURNING *;

-- name: DeleteWorkers :exec
TRUNCATE workers
RESTART IDENTITY;