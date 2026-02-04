#!/usr/bin/bash

# Update a message translation via json

stage=${STAGE:-"stage."}
token=${1?:"Usage: update_translation_json.sh token message_id lang title OR '' text OR ''"}
messageID=${2?:"Usage: update_translation_json.sh token message_id lang title OR '' text OR ''"}
lang=${3?:"Usage: update_translation_json.sh token message_id lang title OR '' text OR ''"}
title=${4?:"Usage: update_translation_json.sh token message_id lang title OR '' text OR ''"}
text=${5?:"Usage: update_translation_json.sh token message_id lang title OR '' text OR ''"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"message\":%MESSAGE%,\"lang\":\"%LANG\",\"text\":\"%TEXT%\",\"title\":\"%TITLE%\"}}' |
	sed -e "s/%TOKEN%/$token/g" -e "s/%MESSAGE%/$messageID/g" -e "s/%LANG%/$lang/g" -e "s/%TEXT%/$text/g" -e "s/%TITLE%/$title/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
sed -e "s/%JSON%/$json/g" ./curl/update_translation_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
