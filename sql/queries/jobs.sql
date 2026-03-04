-- name: GetJob :one
SELECT * FROM jobs
WHERE id = $1;

-- name: GetJobs :many
SELECT * FROM jobs;