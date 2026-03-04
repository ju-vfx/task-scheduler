-- +goose Up
CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT NOT NULL,
    parent_task_id UUID,
    command TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    job_id UUID NOT NULL,
    FOREIGN KEY (job_id)
    REFERENCES jobs(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE tasks;