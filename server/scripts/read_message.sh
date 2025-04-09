#!/usr/bin/bash

# Reads a message

[ -e ".env" ] && source ".env"

stage=${STAGE:-""}
cookie=${COOKIE:-"cookie.txt"}
id=${ID:?"Usage: env ID=<messageId> ./read_message.sh"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%ID%/$id/g" ./curl/read_message.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
