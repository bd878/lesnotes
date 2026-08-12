-- +goose Up
CREATE SCHEMA IF NOT EXISTS messages_stream;

GRANT USAGE ON SCHEMA messages_stream TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT, TRUNCATE ON ALL TABLES IN SCHEMA messages_stream TO lesnotes_admin;

CREATE TABLE IF NOT EXISTS messages_stream.inbox
(
	id          text NOT NULL,
	name        text NOT NULL,
	subject     text NOT NULL,
	data        bytea NOT NULL,
	received_at timestamptz NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS messages_stream.outbox
(
	id           text NOT NULL,
	name         text NOT NULL,
	subject      text NOT NULL,
	data         bytea NOT NULL,
	published_at timestamptz,
	PRIMARY KEY (id)
);

CREATE INDEX messages_unpublished_idx ON messages_stream.outbox (published_at) WHERE published_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS messages_stream.inbox;
DROP TABLE IF EXISTS messages_stream.outbox;
