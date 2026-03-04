-- name: GetJob :one
SELECT * FROM jobs
WHERE id = $1;

-- name: GetJobs :many
SELECT * FROM jobs;

-- name: CreateJob :one
INSERT INTO jobs (
    id, name, status, priority, created_at
) VALUES (
    gen_random_uuid(), $1, $2, $3, NOW()
)
RETURNING *;

-- name: DeleteJobs :exec
DELETE FROM jobs;