#!/usr/bin/bash

# Send a thread message

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
thread_id=${1?"Usage: send_thread_message.sh thread_id"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
sed -e "s/%THREAD_ID%/$thread_id/g" ./curl/send_thread_message.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
