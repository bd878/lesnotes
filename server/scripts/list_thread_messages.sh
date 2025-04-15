#!/usr/bin/bash

# List user thread messages

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
limit=${LIMIT:-"20"}
offset=${OFFSET:-"0"}
thread_id=${1?"Usage: list_thread_message.sh thread_id"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%LIMIT%/$limit/g" \
-e "s/%THREAD_ID%/$thread_id/g" \
-e "s/%OFFSET%/$offset/g" ./curl/list_thread_messages.curl |
curl -b $cookie -s -K -
HERE`
result=`eval "$cmd"`
echo $result
