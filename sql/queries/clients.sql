-- name: GetClient :one
SELECT * FROM clients
WHERE id = $1;

-- name: GetClients :many
SELECT * FROM clients;

-- name: CreateClient :one
INSERT INTO clients (
    host, ip_addr, connected_at, last_seen_at
) VALUES (
    $1, $2, NOW(), NOW()
)
RETURNING *;

-- name: DeleteClients :exec
TRUNCATE clients
RESTART IDENTITY;