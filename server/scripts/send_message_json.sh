#!/usr/bin/bash

# Send a message via json

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}

json='{"text":"abc"}'

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
sed -e "s/%JSON%/$json/g" ./curl/send_message_json.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
