-- +goose Up
CREATE SCHEMA IF NOT EXISTS sessions;

GRANT USAGE ON SCHEMA sessions TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT, TRUNCATE ON ALL TABLES IN SCHEMA sessions TO lesnotes_admin;

CREATE TABLE IF NOT EXISTS sessions.sessions
(
	user_id     bigint       NOT NULL,
	value       VARCHAR(256) NOT NULL,
	expires_at  timestamptz  NOT NULL,
	created_at  timestamptz  NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP SCHEMA IF EXISTS sessions CASCADE;
