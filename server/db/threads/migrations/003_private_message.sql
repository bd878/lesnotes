-- +goose Up
ALTER TABLE threads.threads ADD COLUMN private_message bool NOT NULL DEFAULT true;

-- +goose Down
