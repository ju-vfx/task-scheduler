-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1;

-- name: GetTasks :many
SELECT * FROM tasks;