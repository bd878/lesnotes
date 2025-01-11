PRAGMA user_version=1;
PRAGMA foreigh_keys=ON;
CREATE TABLE IF NOT EXISTS files(
  id                 INTEGER UNIQUE NOT NULL,
  uid                TEXT UNIQUE NOT NULL,
  name               TEXT NOT NULL
);