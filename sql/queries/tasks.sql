-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1;

-- name: GetTasks :many
SELECT * FROM tasks;

-- name: CreateTask :one
INSERT INTO tasks (
    id, name, status, parent_task_id, command, created_at, job_id
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, NOW(), $5
)
RETURNING *;