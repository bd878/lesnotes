-- +goose Up
CREATE SCHEMA IF NOT EXISTS search;

GRANT USAGE ON SCHEMA search TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT, TRUNCATE ON ALL TABLES IN SCHEMA search TO lesnotes_admin;

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

CREATE TRIGGER created_at_search_messages_trgr BEFORE UPDATE ON search.messages FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_search_messages_trgr BEFORE UPDATE ON search.messages FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE IF NOT EXISTS search.threads
(
	id           bigint       UNIQUE NOT NULL,    -- thread id (aka message id)
	user_id      bigint       NOT NULL,
	parent_id    bigint       NOT NULL DEFAULT 0, -- thread id (aka message id)
	name         VARCHAR(256) UNIQUE NOT NULL,
	description  TEXT         NOT NULL DEFAULT '',
	private      bool         NOT NULL DEFAULT true,
	created_at   timestamptz  NOT NULL DEFAULT NOW(),
	updated_at   timestamptz  NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id)
);

CREATE TRIGGER created_at_search_threads_trgr BEFORE UPDATE ON search.threads FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_search_threads_trgr BEFORE UPDATE ON search.threads FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE IF NOT EXISTS search.files
(
	id            bigint        UNIQUE NOT NULL,
	owner_id      bigint        NOT NULL,
	name          VARCHAR(256)  NOT NULL,
	mime          VARCHAR(256)  NOT NULL,
	created_at    timestamptz   NOT NULL DEFAULT NOW(),
	updated_at    timestamptz   NOT NULL DEFAULT NOW(),
	size          int           NOT NULL DEFAULT 0,
	private       bool          NOT NULL DEFAULT true,
	description   TEXT          NOT NULL DEFAULT '',
	PRIMARY KEY(id)
);

CREATE TRIGGER created_at_search_files_trgr BEFORE UPDATE ON search.files FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_search_files_trgr BEFORE UPDATE ON search.files FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE IF NOT EXISTS search.translations
(
	message_id    bigint        NOT NULL,
	user_id       bigint        NOT NULL,
	lang          VARCHAR(8)    NOT NULL,
	text          TEXT          NOT NULL DEFAULT '',
	title         TEXT          NOT NULL DEFAULT '',
	created_at    timestamptz   NOT NULL DEFAULT NOW(),
	updated_at    timestamptz   NOT NULL DEFAULT NOW(),
	PRIMARY KEY(message_id, lang)
);

CREATE TRIGGER created_at_translations_trgr BEFORE UPDATE ON search.translations FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_translations_trgr BEFORE UPDATE ON search.translations FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

-- +goose Down
DROP SCHEMA IF EXISTS search CASCADE;
