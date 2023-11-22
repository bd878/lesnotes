CREATE TABLE IF NOT EXISTS messages(message TEXT, file TEXT);
CREATE TABLE IF NOT EXISTS users(
  id INTEGER PRIMARY KEY,
  name TEXT,
  password TEXT,
  token TEXT,
  expires TEXT
);