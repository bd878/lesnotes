#!/usr/bin/bash

# Sends a message

stage=${STAGE:-""}
cookie=${COOKIE:-"cookie.txt"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" ./curl/send_message.curl |
curl -c $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
