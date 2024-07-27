PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS messages(
  id INTEGER PRIMARY KEY,
  createtime TEXT,
  message TEXT,
  file TEXT,
  file_id TEXT,
  user_id INTEGER
    REFERENCES users(id)
    ON DELETE CASCADE
    NOT NULL,
  log_index INTEGER,
  log_term INTEGER
);
CREATE INDEX IF NOT EXISTS messagesindex ON messages(user_id);
CREATE INDEX IF NOT EXISTS messages_logindex ON messages(log_index, log_term)
  WHERE log_index IS NOT NULL AND log_term IS NOT NULL;
