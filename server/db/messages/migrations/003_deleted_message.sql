-- +goose Up
ALTER TABLE messages.messages ADD COLUMN deleted bool NOT NULL DEFAULT false;
ALTER TABLE messages.translations ADD COLUMN deleted bool NOT NULL DEFAULT false;

-- +goose Down
