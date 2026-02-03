\c lesnotes
CREATE SCHEMA IF NOT EXISTS messages;

-- ORDER IS IMPORTANT!!! FOR SNAPSHOTS
CREATE TABLE IF NOT EXISTS messages.messages
(
	id           bigint        UNIQUE NOT NULL,         -- thread id for child messages
	text         TEXT          NOT NULL,
	file_ids     jsonb         DEFAULT NULL,
	private      bool          NOT NULL DEFAULT true,
	name         VARCHAR(256)  UNIQUE NOT NULL,
	user_id      bigint        NOT NULL,
	-- DELETED: thread_id    bigint        NOT NULL,                -- parent thread id
	created_at   timestamptz   NOT NULL DEFAULT NOW(),
	updated_at   timestamptz   NOT NULL DEFAULT NOW(),
	title        TEXT          NOT NULL DEFAULT '',
	PRIMARY KEY(id)
);

CREATE TRIGGER created_at_messages_trgr BEFORE UPDATE ON messages.messages FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_messages_trgr BEFORE UPDATE ON messages.messages FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

-- TODO: messages.comments, messages.reactions

CREATE TABLE IF NOT EXISTS messages.files
(
	file_id       bigint       NOT NULL,
	message_id    bigint       NOT NULL,
	user_id       bigint       NOT NULL,
	PRIMARY KEY(file_id, message_id)
);

-- TODO: add foreign key constraint files -> messages

CREATE TABLE IF NOT EXISTS messages.translations
(
	message_id    bigint       NOT NULL,
	lang          VARCHAR(8)   NOT NULL,
	text          TEXT         NOT NULL DEFAULT '',
	title         TEXT         NOT NULL DEFAULT '',
	created_at    timestamptz   NOT NULL DEFAULT NOW(),
	updated_at    timestamptz   NOT NULL DEFAULT NOW(),
	PRIMARY KEY(message_id, lang)
);

CREATE TRIGGER created_at_translations_trgr BEFORE UPDATE ON messages.translations FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_translations_trgr BEFORE UPDATE ON messages.translations FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

GRANT USAGE ON SCHEMA messages TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT, TRUNCATE ON ALL TABLES IN SCHEMA messages TO lesnotes_admin;
