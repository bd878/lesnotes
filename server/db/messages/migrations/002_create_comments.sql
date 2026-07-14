-- +goose Up
CREATE TABLE IF NOT EXISTS messages.comments
(
	id           bigint        UNIQUE NOT NULL,
	message_id   bigint        NOT NULL,
	user_id      bigint        NOT NULL,
	text         TEXT          NOT NULL,
	metadata     bytea         DEFAULT NULL,
	created_at   timestamptz   NOT NULL DEFAULT NOW(),
	updated_at   timestamptz   NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id, user_id)
);

-- +goose Down
DROP TABLE IF EXISTS messages.comments;