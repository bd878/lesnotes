#!/usr/bin/bash

# Deletes a messages

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
id=${1?"Usage: delete_messages.sh ...message_id"}

argv=( "$@" )
ids="%5B${argv[0]}"
for arg in "${argv[@]:1}"; do
	ids="${ids}%2C$arg"
done
ids="${ids}%5D"

echo -n "$ids"

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
sed -e "s/%IDS%/$ids/g" ./curl/delete_messages.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
