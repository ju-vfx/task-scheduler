-- name: GetWorker :one
SELECT * FROM workers
WHERE id = $1;

-- name: GetWorkers :many
SELECT * FROM workers;

-- name: CreateWorker :one
INSERT INTO workers (
    id, host, port, connected_at, last_seen_at
) VALUES (
    gen_random_uuid(), $1, $2, NOW(), NOW()
)
RETURNING *;

-- name: UpdateLastSeen :exec
UPDATE workers
SET last_seen_at = NOW()
WHERE id = $1;

-- name: DeleteWorkers :exec
DELETE FROM workers;