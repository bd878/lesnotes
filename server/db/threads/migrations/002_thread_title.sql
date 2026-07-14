-- +goose Up
ALTER TABLE threads.threads ADD COLUMN title text NOT NULL DEFAULT '';

-- +goose Down
