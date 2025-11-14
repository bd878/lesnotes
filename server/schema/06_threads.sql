\c lesnotes
CREATE SCHEMA IF NOT EXISTS threads;

CREATE TABLE IF NOT EXISTS threads.threads
(
	id           bigint        UNIQUE NOT NULL,   -- message id
	name         VARCHAR(256)  UNIQUE NOT NULL,
	private      bool          NOT NULL DEFAULT true,
	user_id      bigint        NOT NULL,
	parent_id    bigint        NOT NULL,          -- parent thread id (aka parent message id)
	created_at   timestamptz   NOT NULL DEFAULT NOW(),
	updated_at   timestamptz   NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id)
);

CREATE TRIGGER created_at_threads_trgr BEFORE UPDATE ON threads.threads FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_threads_trgr BEFORE UPDATE ON threads.threads FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

GRANT USAGE ON SCHEMA threads TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA threads TO lesnotes_admin;
