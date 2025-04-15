#!/usr/bin/bash

# Send a public message

stage=${STAGE:-"stage."}
cookie=${COOKIE:-""}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" ./curl/send_public_message.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
