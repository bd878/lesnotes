-- +goose Up
CREATE SCHEMA IF NOT EXISTS users_stream;

GRANT USAGE ON SCHEMA users_stream TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT, TRUNCATE ON ALL TABLES IN SCHEMA users_stream TO lesnotes_admin;

CREATE TABLE IF NOT EXISTS users_stream.inbox
(
	id          text NOT NULL,
	name        text NOT NULL,
	subject     text NOT NULL,
	data        bytea NOT NULL,
	received_at timestamptz NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS users_stream.outbox
(
	id           text NOT NULL,
	name         text NOT NULL,
	subject      text NOT NULL,
	data         bytea NOT NULL,
	published_at timestamptz,
	PRIMARY KEY (id)
);

CREATE INDEX users_unpublished_idx ON users_stream.outbox (published_at) WHERE published_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS users_stream.inbox;
DROP TABLE IF EXISTS users_stream.outbox;
