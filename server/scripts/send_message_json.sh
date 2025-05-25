#!/usr/bin/bash

# Send a message via json

stage=${STAGE:-"stage."}
token=${1?:"Usage: send_message_json.sh token"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"text\":\"abc\"}}' | sed -e "s/%TOKEN%/$token/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
sed -e "s/%JSON%/$json/g" ./curl/send_message_json.curl |
curl -v -K -
HERE`
result=`eval "$cmd"`
echo $result
