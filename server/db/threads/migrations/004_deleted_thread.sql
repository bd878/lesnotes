-- +goose Up
ALTER TABLE threads.threads ADD COLUMN deleted bool NOT NULL DEFAULT false;

-- +goose Down
