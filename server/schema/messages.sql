PRAGMA user_version=1;
PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS messages(
  id                  INTEGER UNIQUE NOT NULL,
  create_utc_nano     INTEGER NOT NULL,
  update_utc_nano     INTEGER NOT NULL,
  text                TEXT,
  file_id             INTEGER REFERENCES files(id) ON DELETE CASCADE,
  user_id             INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL
);
CREATE INDEX IF NOT EXISTS messages_index ON messages(user_id);

