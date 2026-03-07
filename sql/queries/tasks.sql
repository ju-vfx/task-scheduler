-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1;

-- name: GetTasks :many
SELECT * FROM tasks
ORDER BY created_at ASC;

-- name: GetTasksByJobId :many
SELECT * FROM tasks
WHERE job_id = $1
ORDER BY created_at ASC;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET status = $2
WHERE id = $1
RETURNING *;

-- name: CreateTask :one
INSERT INTO tasks (
    id, name, status, command, created_at, job_id
) VALUES (
    gen_random_uuid(), $1, $2, $3, NOW(), $4
)
RETURNING *;