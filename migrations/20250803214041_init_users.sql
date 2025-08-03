-- +goose Up
CREATE TABLE users (
                       id BIGINT PRIMARY KEY,
                       first_name TEXT,
                       last_name TEXT,
                       username TEXT,
                       language TEXT,
                       created_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS users;