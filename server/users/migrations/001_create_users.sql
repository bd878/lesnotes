-- +goose Up
CREATE SCHEMA IF NOT EXISTS users;

GRANT USAGE ON SCHEMA users TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT, TRUNCATE ON ALL TABLES IN SCHEMA users TO lesnotes_admin;

CREATE TABLE IF NOT EXISTS users.users
(
	id          bigint       UNIQUE NOT NULL,
	login       TEXT         UNIQUE NOT NULL,
	salt        TEXT         NOT NULL,
	metadata    bytea        DEFAULT NULL,
	blocked     bool         NOT NULL DEFAULT false,
	created_at  timestamptz  NOT NULL DEFAULT NOW(),
	updated_at  timestamptz  NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id)
);

-- +goose Down
DROP SCHEMA IF EXISTS users CASCADE;
