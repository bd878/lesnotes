ALTER TABLE messages ADD COLUMN file_id TEXT;
ALTER TABLE messages ADD COLUMN log_index INTEGER;
ALTER TABLE messages ADD COLUMN log_term INTEGER;
CREATE INDEX IF NOT EXISTS messages_logindex ON messages(log_index, log_term)
  WHERE log_index IS NOT NULL AND log_term IS NOT NULL;