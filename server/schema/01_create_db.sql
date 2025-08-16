CREATE DATABASE lesnotes;

GRANT CONNECT ON DATABASE lesnotes TO lesnotes_admin;
ALTER ROLE lesnotes_admin SET search_path TO lesnotes, "$user", public;

\c lesnotes

CREATE OR REPLACE FUNCTION created_at_trigger()
RETURNS TRIGGER AS $$
BEGIN
	NEW.created_at := OLD.created_at;
	RETURN NEW;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION updated_at_trigger()
RETURNS TRIGGER AS $$
BEGIN
	IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
		NEW.updated_at = NOW();
		RETURN NEW;
	ELSE
		RETURN OLD;
	END IF;
END;
$$ language 'plpgsql';