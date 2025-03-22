#!/bin/bash

# Sets up sqlite database,
# runs all files migrations

DB_FILE=${1?"Usage: $0 sqlite_file.db"}

echo -n "Setup $DB_FILE? (y/N) "
read agree
case "$agree" in
	[yY]) echo "Running migration on $DB_FILE...";;
	
	*) echo "No changes made, exit."
		 exit 0
		 ;;
esac

sqlite3 $DB_FILE < ./schema/files.sql

echo "done."

exit 0;
