-- +goose Up
CREATE SCHEMA IF NOT EXISTS files_stream;

GRANT USAGE ON SCHEMA files_stream TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT, TRUNCATE ON ALL TABLES IN SCHEMA files_stream TO lesnotes_admin;

CREATE TABLE IF NOT EXISTS files_stream.inbox
(
	id          text NOT NULL,
	name        text NOT NULL,
	subject     text NOT NULL,
	data        bytea NOT NULL,
	received_at timestamptz NOT NULL,
	PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE IF EXISTS files_stream.inbox;
