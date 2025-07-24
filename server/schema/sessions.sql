PRAGMA user_version=1;
PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS sessions(
	user_id          INTEGER UNIQUE NOT NULL,
	value            TEXT NOT NULL,
	expires_utc_nano INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS sessions_index ON sessions(user_id);
