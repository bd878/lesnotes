-- +goose Up
CREATE SCHEMA IF NOT EXISTS users;

GRANT USAGE ON SCHEMA users TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT, TRUNCATE ON ALL TABLES IN SCHEMA users TO lesnotes_admin;

CREATE TABLE IF NOT EXISTS users.users
(
	id          bigint       UNIQUE NOT NULL,
	login       TEXT         UNIQUE NOT NULL,
	salt        TEXT         NOT NULL,
	created_at  timestamptz  NOT NULL DEFAULT NOW(),
	updated_at  timestamptz  NOT NULL DEFAULT NOW(),
	metadata    bytea        DEFAULT NULL,
	blocked     bool         NOT NULL DEFAULT false,
	PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS users_id ON users.users(id);

CREATE TRIGGER created_at_users_trgr BEFORE UPDATE ON users.users FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_users_trgr BEFORE UPDATE ON users.users FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

-- +goose Down
DROP SCHEMA IF EXISTS users CASCADE;
