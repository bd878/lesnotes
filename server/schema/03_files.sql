\с lesnotes
CREATE SCHEMA IF NOT EXISTS files;

CREATE TABLE IF NOT EXISTS files.files
(
	id            bigint        UNIQUE NOT NULL,
	user_id       bigint        NOT NULL,
	name          VARCHAR(256)  UNIQUE NOT NULL,
	private       bool          NOT NULL DEFAULT true,
	created_at    timestamptz   NOT NULL DEFAULT NOW(),
	updated_at    timestamptz   NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS files_id ON files.files(id);

CREATE TRIGGER created_at_files_trgr BEFORE UPDATE ON files.files FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_files_trgr BEFORE UPDATE ON files.files FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

GRANT USAGE ON SCHEMA files TO lesnotes_admin;
GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA files TO lesnotes_admin;