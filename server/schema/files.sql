PRAGMA user_version=1;
PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS files(
	id                 INTEGER UNIQUE NOT NULL,
	user_id            INTEGER NOT NULL,
	name               TEXT NOT NULL,
	create_utc_nano    INTEGER NOT NULL
);