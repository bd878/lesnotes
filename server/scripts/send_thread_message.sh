#!/usr/bin/bash

# Send a thread message

[ -e ".env" ] && source ".env"

stage=${STAGE:-""}
cookie=${COOKIE:-"cookie.txt"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" ./curl/send_message.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
