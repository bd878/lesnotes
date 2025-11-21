#!/usr/bin/bash

# Read thread json request

stage=${STAGE:-"stage."}
token=${1?"Usage: read_thread_json.sh token user_id id"}
user=${2?"Usage: read_thread_json.sh token user_id id"}
id=${3?"Usage: read_thread_json.sh token user_id id"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"id\":%ID%,\"user_id\":%USER%}}' |
	sed -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g" -e "s/%USER%/$user/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/read_thread_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
