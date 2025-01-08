#!/bin/bash

# Sets up sqlite database,
# runs all messages migrations

DB_FILE=${1?"Usage: $0 sqlite_file.db"}

echo -n "Setup $DB_FILE? (y/N) "
read agree
case "$agree" in
  [yY]) echo "Running migration on $DB_FILE...";;
  
  *) echo "No changes made, exit."
     exit 0
     ;;
esac

sqlite3 $DB_FILE < ./schema/messages.sql
sqlite3 $DB_FILE < ./schema/messages_add_log_columns.sql
sqlite3 $DB_FILE < ./schema/messages_add_update_utc_nano_column.sql

echo "done."

exit 0;
