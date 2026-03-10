-- +goose Up
CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    status INT NOT NULL,
    priority INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP,
    cancelled_at TIMESTAMP
);

-- +goose Down
DROP TABLE jobs;