PRAGMA user_version=1;
PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS users(
  id                    INTEGER UNIQUE NOT NULL,
  name                  TEXT,
  password              TEXT,
  token                 TEXT,
  expires               TEXT
);