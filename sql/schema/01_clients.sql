-- +goose Up
CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    host TEXT NOT NULL,
    ip_addr TEXT NOT NULL,
    connected_at TIMESTAMP NOT NULL,
    last_seen_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE clients;