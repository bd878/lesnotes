#!/usr/bin/bash

# Publish a given private thread via json

stage=${STAGE:-"stage."}
token=${1?"Usage: public_thread_json.sh token thread_id"}
thread=${2?"Usage: public_thread_json.sh token thread_id"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"id\":%ID%}}' | sed -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/public_thread_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
