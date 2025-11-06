\c lesnotes
CREATE SCHEMA IF NOT EXISTS search;

CREATE TABLE IF NOT EXISTS search.messages
(
	id         bigint       UNIQUE NOT NULL,    -- message id
	user_id    bigint       NOT NULL,
	name       VARCHAR(256) UNIQUE NOT NULL,
	text       TEXT         NOT NULL DEFAULT '',
	title      TEXT         NOT NULL DEFAULT '',
	private    bool         NOT NULL DEFAULT true,
	created_at timestamptz  NOT NULL DEFAULT NOW(),
	updated_at timestamptz  NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS search_messages_text ON search.messages(text);

CREATE TRIGGER created_at_search_messages_trgr BEFORE UPDATE ON search.messages1 FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_search_messages_trgr BEFORE UPDATE ON search.messages1 FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE IF NOT EXISTS search.files
(
	id          bigint       UNIQUE NOT NULL,
	owner_id    bigint       NOT NULL,
	name        VARCHAR(256) NOT NULL,
	mime        VARCHAR(256) NOT NULL,
	created_at    timestamptz   NOT NULL DEFAULT NOW(),
	updated_at    timestamptz   NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS search_files_name ON search.files(name);

CREATE TRIGGER created_at_search_files_trgr BEFORE UPDATE ON search.files FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_search_files_trgr BEFORE UPDATE ON search.files FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

GRANT USAGE ON SCHEMA search TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA search TO lesnotes_admin;
