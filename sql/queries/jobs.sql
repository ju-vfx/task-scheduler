-- name: GetJob :one
SELECT * FROM jobs
WHERE id = $1;

-- name: GetJobs :many
SELECT * FROM jobs
ORDER BY created_at DESC;

-- name: GetWaitingJobs :many
SELECT * FROM jobs
WHERE finished_at IS NULL
AND cancelled_at IS NULL
ORDER BY priority DESC, created_at ASC;

-- name: UpdateJobStatus :exec
UPDATE jobs
SET status = $2, finished_at = $3, cancelled_at = $4
WHERE id = $1;

-- name: CreateJob :one
INSERT INTO jobs (
    id, name, status, priority, created_at
) VALUES (
    gen_random_uuid(), $1, $2, $3, NOW()
)
RETURNING *;

-- name: DeleteJobs :exec
DELETE FROM jobs;