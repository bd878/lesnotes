#!/usr/bin/bash

# Update a message.
# Move it to other thread

usage="Usage: move_message.sh message_id thread_id"

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
id=${1?"$usage"}
thread_id=${2:-"-1"}

printf "%s %s\n" $id $thread_id

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%ID%/$id/g" \
-e "s/%THREAD_ID%/$thread_id/g" ./curl/move_message.curl |
curl  --trace-ascii /dev/stdout -b $cookie -K -
HERE`
result=`eval "$cmd"`
echo $result
