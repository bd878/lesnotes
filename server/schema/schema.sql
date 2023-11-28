PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS users(
  id INTEGER PRIMARY KEY,
  name TEXT,
  password TEXT,
  token TEXT,
  expires TEXT
);
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
CREATE INDEX messagesindex ON messages(user_id);
