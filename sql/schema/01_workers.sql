-- +goose Up
CREATE TABLE workers (
    id UUID PRIMARY KEY,
    host TEXT NOT NULL,
    port TEXT NOT NULL,
    connected_at TIMESTAMP NOT NULL,
    last_seen_at TIMESTAMP NOT NULL,
    status TEXT
);

-- +goose Down
DROP TABLE workers;