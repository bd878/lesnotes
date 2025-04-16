#!/usr/bin/bash

# List user messages

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
limit=${1?"Usage: list_messages.sh limit offset"}
offset=${2?"Usage: list_messages.sh limit offset"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%LIMIT%/$limit/g" \
-e "s/%OFFSET%/$offset/g" ./curl/list_messages.curl |
curl -b $cookie -s -K -
HERE`
result=`eval "$cmd"`
echo $result
