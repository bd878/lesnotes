#!/usr/bin/bash

# Update a comment via json

stage=${STAGE:-"stage."}
token=${1?:"Usage: update_comment_json.sh token id text"}
id=${2?:"Usage: update_comment_json.sh token id text"}
text=${3?:"Usage: update_comment_json.sh token id text"}

json=$(echo -n '
{
  "token": "%TOKEN%",
  "req": {
    "id": %ID%,
    "text": "%TEXT%"
  }
}
' | tr -d ' ' | tr -d '\n' | sed -e 's/"/\\\"/g' -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g" -e "s/%TEXT%/$text/g" )

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" -e "s/%JSON%/$(echo -n $json)/g" ./curl/update_comment_json.curl | curl -K -
HERE`
result=`eval "$cmd"`
echo $result
