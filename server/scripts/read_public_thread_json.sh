#!/usr/bin/bash

# Read public thread json request

stage=${STAGE:-"stage."}
user=${1?"Usage: read_thread_json.sh user_id AND id OR name"}
id=${2?"Usage: read_thread_json.sh user_id AND id OR name"}
name=${3?"Usage: read_thread_json.sh user_id AND id OR name"}

json=$(echo -n '{\"token\":\"""\",\"req\":{\"id\":%ID%,\"user_id\":%USER%,\"name\":\"%NAME%\"}}' |
	sed -e "s/%ID%/$id/g" -e "s/%USER%/$user/g" -e "s/%NAME%/$name/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/read_thread_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
