#!/usr/bin/bash

# Reads a message

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
id=${1:?"Usage: ./read_message.sh message_id"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%ID%/$id/g" ./curl/read_message.curl |
curl -b $cookie -K -
HERE`
result=`eval "$cmd"`
echo $result
