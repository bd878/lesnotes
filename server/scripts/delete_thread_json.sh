#!/usr/bin/bash

# Delete thread json request

stage=${STAGE:-"stage."}
token=${1?"Usage: delete_thread_json.sh token id"}
id=${2?"Usage: delete_thread_json.sh token id"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"id\":%ID%}}' |
	sed -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/delete_thread_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
