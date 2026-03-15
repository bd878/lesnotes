#!/usr/bin/bash

# Send a comment via json

stage=${STAGE:-"stage."}
token=${1?:"Usage: send_comment_json.sh token message_id text"}
messageID=${2?:"Usage: send_comment_json.sh token message_id text"}
text=${3?:"Usage: send_comment_json.sh token message_id text"}

json=$(echo -n '
{
  "token": "%TOKEN%",
  "req": {
    "message": %MESSAGE%,
    "text": "%TEXT%"
  }
}
' | tr -d ' ' | tr -d '\n' | sed -e 's/"/\\\"/g' -e "s/%TOKEN%/$token/g" -e "s/%MESSAGE%/$messageID/g" -e "s/%TEXT%/$text/g" )

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" -e "s/%JSON%/$(echo -n $json)/g" ./curl/send_comment_json.curl | curl -K -
HERE`
result=`eval "$cmd"`
echo $result
