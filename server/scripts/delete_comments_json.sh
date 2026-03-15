#!/usr/bin/bash

# Delete a comment via json

stage=${STAGE:-"stage."}
token=${1?:"Usage: delete_comments_json.sh token id "}
id=${2?:"Usage: delete_comments_json.sh token id "}

json=$(echo -n '
{
  "token": "%TOKEN%",
  "req": {
    "id": %ID%
  }
}
' | tr -d ' ' | tr -d '\n' | sed -e 's/"/\\\"/g' -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" -e "s/%JSON%/$(echo -n $json)/g" ./curl/delete_comments_json.curl | curl -K -
HERE`
result=`eval "$cmd"`
echo $result
