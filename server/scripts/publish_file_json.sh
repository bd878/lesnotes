#!/usr/bin/bash

# Publish a given private file via json

stage=${STAGE:-"stage."}
token=${1?"Usage: publish_file_json.sh token file_id"}
id=${2?"Usage: publish_file_json.sh token file_id"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"id\":%ID%}}' | sed -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/publish_file_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
