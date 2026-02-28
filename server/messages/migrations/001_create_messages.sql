-- +goose Up
CREATE SCHEMA IF NOT EXISTS messages;

GRANT USAGE ON SCHEMA messages TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT, TRUNCATE ON ALL TABLES IN SCHEMA messages TO lesnotes_admin;

CREATE TABLE IF NOT EXISTS messages.messages
(
	id           bigint        UNIQUE NOT NULL,
	text         TEXT          NOT NULL,
	private      bool          NOT NULL DEFAULT true,
	name         VARCHAR(256)  UNIQUE NOT NULL,
	user_id      bigint        NOT NULL,
	created_at   timestamptz   NOT NULL DEFAULT NOW(),
	updated_at   timestamptz   NOT NULL DEFAULT NOW(),
	title        TEXT          NOT NULL DEFAULT '',
	PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS messages.files
(
	file_id       bigint       NOT NULL,
	message_id    bigint       NOT NULL,
	user_id       bigint       NOT NULL,
	PRIMARY KEY(file_id, message_id)
);

CREATE TABLE IF NOT EXISTS messages.translations
(
	message_id    bigint        NOT NULL,
	lang          VARCHAR(8)    NOT NULL,
	text          TEXT          NOT NULL DEFAULT '',
	title         TEXT          NOT NULL DEFAULT '',
	created_at    timestamptz   NOT NULL DEFAULT NOW(),
	updated_at    timestamptz   NOT NULL DEFAULT NOW(),
	PRIMARY KEY(message_id, lang)
);

-- +goose Down
DROP SCHEMA IF EXISTS messages CASCADE;
