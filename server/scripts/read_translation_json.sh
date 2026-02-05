#!/usr/bin/bash

# Read a message translation via json

stage=${STAGE:-"stage."}
token=${1?:"Usage: read_translation_json.sh token message_id lang"}
messageID=${2?:"Usage: read_translation_json.sh token message_id lang"}
lang=${3?:"Usage: read_translation_json.sh token message_id lang"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"message\":%MESSAGE%,\"lang\":\"%LANG\"}}' |
	sed -e "s/%TOKEN%/$token/g" -e "s/%MESSAGE%/$messageID/g" -e "s/%LANG%/$lang/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
sed -e "s/%JSON%/$json/g" ./curl/read_translation_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
