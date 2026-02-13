-- +goose Up
CREATE SCHEMA IF NOT EXISTS sessions;

CREATE TABLE IF NOT EXISTS sessions.sessions
(
	user_id     bigint       NOT NULL,
	value       VARCHAR(256) NOT NULL,
	expires_at  timestamptz  NOT NULL,
	created_at  timestamptz  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER created_at_sessions_trgr BEFORE UPDATE ON sessions.sessions FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();

GRANT USAGE ON SCHEMA sessions TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT, TRUNCATE ON ALL TABLES IN SCHEMA sessions TO lesnotes_admin;

-- +goose Down
DROP SCHEMA IF EXISTS sessions CASCADE;
