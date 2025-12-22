#!/usr/bin/bash

# Private a given public thread via json

stage=${STAGE:-"stage."}
token=${1?"Usage: private_thread_json.sh token thread_id"}
id=${2?"Usage: private_thread_json.sh token thread_id"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"id\":%ID%}}' | sed -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/private_thread_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
