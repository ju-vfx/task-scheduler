-- name: GetJob :one
SELECT * FROM jobs
WHERE id = $1;

-- name: GetJobs :many
SELECT * FROM jobs;

-- name: GetWaitingJobs :many
SELECT * FROM jobs
WHERE finished_at IS NULL
AND cancelled_at IS NULL
ORDER BY priority, created_at DESC;

-- name: UpdateJobFinished :exec
UPDATE jobs
SET status = $2, finished_at = NOW()
WHERE id = $1;

-- name: UpdateJobError :exec
UPDATE jobs
SET status = $2, cancelled_at = NOW()
WHERE id = $1;

-- name: UpdateJobRunning :exec
UPDATE jobs
SET status = $2
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