\c lesnotes
CREATE SCHEMA IF NOT EXISTS users;

CREATE TYPE users.theme AS ENUM ('dark', 'light');

CREATE TABLE IF NOT EXISTS users.users
(
	id          bigint       UNIQUE NOT NULL,
	login       TEXT         UNIQUE NOT NULL,
	salt        TEXT         NOT NULL,
	theme       users.theme  NOT NULL DEFAULT 'light',
	lang        VARCHAR(3)   NOT NULL,
	font_size   int          NOT NULL,
	created_at  timestamptz  NOT NULL DEFAULT NOW(),
	updated_at  timestamptz  NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS users_id ON users.users(id);

CREATE TRIGGER created_at_users_trgr BEFORE UPDATE ON users.users FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_users_trgr BEFORE UPDATE ON users.users FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

GRANT USAGE ON SCHEMA users TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA users TO lesnotes_admin;