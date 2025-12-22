#!/usr/bin/bash

# Read thread json request

stage=${STAGE:-"stage."}
token=${1?"Usage: read_thread_json.sh token OR user_id id OR name"}
user=${2?"Usage: read_thread_json.sh token OR user_id id OR name"}
id=${3?"Usage: read_thread_json.sh token OR user_id id OR name"}
name=${4?"Usage: read_thread_json.sh token OR user_id id OR name"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"id\":%ID%,\"user_id\":%USER%,\"name\":\"%NAME%\"}}' |
	sed -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g" -e "s/%USER%/$user/g" -e "s/%NAME%/$name/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/read_thread_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
