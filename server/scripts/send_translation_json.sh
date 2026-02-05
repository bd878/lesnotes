#!/usr/bin/bash

# Send a message translation via json

stage=${STAGE:-"stage."}
token=${1?:"Usage: send_translation_json.sh token message_id lang title text"}
messageID=${2?:"Usage: send_translation_json.sh token message_id lang title text"}
lang=${3?:"Usage: send_translation_json.sh token message_id lang title text"}
title=${4?:"Usage: send_translation_json.sh token message_id lang title text"}
text=${5?:"Usage: send_translation_json.sh token message_id lang title text"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"message\":%MESSAGE%,\"lang\":\"%LANG%\",\"text\":\"%TEXT%\",\"title\":\"%TITLE%\"}}' |
	sed -e "s/%TOKEN%/$token/g" -e "s/%MESSAGE%/$messageID/g" -e "s/%LANG%/$lang/g" -e "s/%TEXT%/$text/g" -e "s/%TITLE%/$title/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
sed -e "s/%JSON%/$json/g" ./curl/send_translation_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
