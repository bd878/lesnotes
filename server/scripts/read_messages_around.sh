#!/usr/bin/bash

# List messages around one message in a thread

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
thread_id=${1?"Usage: read_messages_around.sh thread_id id limit"}
id=${2?"Usage: read_messages_around.sh thread_id id limit"}
limit=${3?"Usage: read_messages_arounds.sh thread_id id limit"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%LIMIT%/$limit/g" \
-e "s/%THREAD_ID%/$thread_id/g" \
-e "s/%ID%/$id/g" ./curl/read_messages_around.curl |
curl -b $cookie -s -K -
HERE`
result=`eval "$cmd"`
echo $result
