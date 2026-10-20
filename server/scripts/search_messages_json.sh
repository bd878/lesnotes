#!/usr/bin/bash

# Search messages by query
# via json api

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
token=${1?"Usage: search_messages_json.sh token query name"}
query=${2?"Usage: search_messages_json.sh token query name"}
name=${3?"Usage: search_messages_json.sh token query name"}

json=$(echo -n '
{
  "token": "%TOKEN%",
  "req": {
    "query": "%QUERY%",
    "name": "%NAME%"
  }
}' | tr -d ' ' | tr -d '\n' | sed -e 's/"/\\\"/g' -e "s/%TOKEN%/$token/g" -e "s/%QUERY%/$query/g" \
        -e "s/%NAME%/$name/g" )

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$(echo -n $json)/g" ./curl/search_messages_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
