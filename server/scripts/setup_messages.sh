#!/bin/bash

# Sets up sqlite database,
# runs all messages migrations

DB_FILE=${1?"Usage: $0 sqlite_file.db"}

run_migrations() {
	local file=${1?"Usage: run_migrations db_file.sql"}
	for f in `ls ./migrations/*.sql`
	do
		printf "%s\n" $f
		sqlite3 $file < $f
	done
}

echo -n "Setup $DB_FILE? (y/N) "
read agree
case "$agree" in
	[yY])
		echo "Running migration on $DB_FILE..."
		run_migrations $DB_FILE
		;;

	*)
		echo "No changes made, exit."
		exit 0
		;;
esac

sqlite3 $DB_FILE < ./schema/messages.sql

echo "done."

exit 0;
