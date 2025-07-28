PRAGMA user_version=1;
PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS sessions(
	user_id          INTEGER NOT NULL,
	value            TEXT UNIQUE NOT NULL,
	expires_utc_nano INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS sessions_index ON sessions(user_id);
