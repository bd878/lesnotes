#!/usr/bin/bash

# List user messages via json request

stage=${STAGE:-"stage."}
token=${1?"Usage: list_messages_json.sh token"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"limit\":10,\"offset\":0}}' | sed -e "s/%TOKEN%/$token/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/list_messages_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
