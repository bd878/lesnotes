-- +goose Up
CREATE SCHEMA IF NOT EXISTS files;

GRANT CREATE ON DATABASE lesnotes TO lesnotes_admin;
GRANT CREATE ON SCHEMA public TO lesnotes_admin;
GRANT USAGE ON SCHEMA files TO lesnotes_admin;
GRANT CREATE ON SCHEMA files TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA files TO lesnotes_admin;

CREATE TABLE IF NOT EXISTS files.files
(
	id            bigint        UNIQUE NOT NULL,
	owner_id      bigint        NOT NULL,
	name          VARCHAR(256)  NOT NULL,
	private       bool          NOT NULL DEFAULT true,
	oid           int           UNIQUE DEFAULT NULL, -- large object id
	mime          VARCHAR(256)  NOT NULL,
	size          int           NOT NULL,
	created_at    timestamptz   NOT NULL DEFAULT NOW(),
	updated_at    timestamptz   NOT NULL DEFAULT NOW(),
	description   text          NOT NULL DEFAULT '',
	PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS files.messages
(
	file_id       bigint       NOT NULL,
	message_id    bigint       NOT NULL,
	user_id       bigint       NOT NULL,
	PRIMARY KEY(file_id, message_id)
);

-- +goose Down
