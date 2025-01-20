PRAGMA user_version=1;
PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS messages(
  id                  INTEGER UNIQUE NOT NULL,
  create_utc_nano     INTEGER NOT NULL,
  update_utc_nano     INTEGER NOT NULL,
  text                TEXT,
  file_id             INTEGER DEFAULT NULL,
  user_id             INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS messages_index ON messages(user_id);
