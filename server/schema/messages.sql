PRAGMA user_version=1;
PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS messages(
  id INTEGER PRIMARY KEY,
  createtime TEXT,
  message TEXT,
  file TEXT,
  user_id INTEGER
    REFERENCES users(id)
    ON DELETE CASCADE
    NOT NULL
);
CREATE INDEX IF NOT EXISTS messagesindex ON messages(user_id);
