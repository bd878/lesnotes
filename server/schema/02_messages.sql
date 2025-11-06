\c lesnotes
CREATE SCHEMA IF NOT EXISTS messages;

CREATE TABLE IF NOT EXISTS messages.messages
(
	id           bigint        UNIQUE NOT NULL,         -- thread id for child messages
	text         TEXT          NOT NULL,
	title        TEXT          NOT NULL DEFAULT '',
	file_ids     jsonb         DEFAULT NULL,
	private      bool          NOT NULL DEFAULT true,
	name         VARCHAR(256)  UNIQUE NOT NULL,
	user_id      bigint        NOT NULL,
	thread_id    bigint        NOT NULL,                -- parent message id
	created_at   timestamptz   NOT NULL DEFAULT NOW(),
	updated_at   timestamptz   NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id)
);

CREATE TRIGGER created_at_messages_trgr BEFORE UPDATE ON messages.messages FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_messages_trgr BEFORE UPDATE ON messages.messages FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

GRANT USAGE ON SCHEMA messages TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT, TRUNCATE ON ALL TABLES IN SCHEMA messages TO lesnotes_admin;
