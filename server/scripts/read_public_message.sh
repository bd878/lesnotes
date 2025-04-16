#!/usr/bin/bash

# Reads a public message

stage=${STAGE:-"stage."}
cookie=${COOKIE:-""}
id=${1:?"Usage: ./read_message.sh message_id"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%ID%/$id/g" ./curl/read_public_message.curl |
curl -v -K -
HERE`
result=`eval "$cmd"`
echo $result
