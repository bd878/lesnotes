#!/usr/bin/bash

# Reads a messages tree

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
token=${1:?"Usage: ./read_tree_json.sh token message_id limit offset"}
id=${2:?"Usage: ./read_tree_json.sh token message_id limit offset"}
limit=${3:?"Usage: ./read_tree_json.sh token message_id limit offset"}
offset=${4:?"Usage: ./read_tree_json.sh token message_id limit offset"}

json=$(echo -n '
{
  "token": "%TOKEN%",
  "req": {
    "root": %ID%,
    "limit": %LIMIT%,
    "offset": %OFFSET%,
    "leaves": [
        {
            "id": 1615442471,
            "limit": 10,
            "offset": 0
        }
    ]
  }
}' | tr -d ' ' | tr -d '\n' | sed -e 's/"/\\\"/g' -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g" \
        -e "s/%LIMIT%/$limit/g" -e "s/%OFFSET%/$offset/g" )

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" -e "s/%JSON%/$(echo -n $json)/g" ./curl/read_tree_json.curl | curl -K -
HERE`
# echo $cmd
result=`eval "$cmd"`
echo $result
