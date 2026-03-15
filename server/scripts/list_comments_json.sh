#!/usr/bin/bash

# List message comments via json

stage=${STAGE:-"stage."}
token=${1?:"Usage: list_comments_json.sh token message_id limit offset"}
messageID=${2?:"Usage: list_comments_json.sh token message_id limit offset"}
limit=${3?:"Usage: list_comments_json.sh token message_id limit offset"}
offset=${4?:"Usage: list_comments_json.sh token message_id limit offset"}

json=$(echo -n '
{
  "token": "%TOKEN%",
  "req": {
    "message": %MESSAGE%,
    "limit": %LIMIT%,
    "offset": %OFFSET%
  }
}
' | tr -d ' ' | tr -d '\n' | sed -e 's/"/\\\"/g' -e "s/%TOKEN%/$token/g" \
  -e "s/%MESSAGE%/$messageID/g" -e "s/%LIMIT%/$limit/g" -e "s/%OFFSET%/$offset/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" -e "s/%JSON%/$(echo -n $json)/g" ./curl/list_comments_json.curl | curl -K -
HERE`
result=`eval "$cmd"`
echo $result
