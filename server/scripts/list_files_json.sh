#!/usr/bin/bash

# List user files via json request

stage=${STAGE:-"stage."}
token=${1?"Usage: list_files_json.sh token OR '' user_id OR 0 limit offset"}
user=${2?"Usage: list_files_json.sh token OR '' user_id OR 0 limit offset"}
limit=${3?"Usage: list_files_json.sh token OR '' user_id OR 0 limit offset"}
offset=${4?"Usage: list_files_json.sh token OR '' user_id OR 0 limit offset"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"limit\":%LIMIT%,\"offset\":%OFFSET%,\"user\":%USER%}}' | 
	sed -e "s/%TOKEN%/$token/g" -e "s/%LIMIT%/$limit/g" -e "s/%OFFSET%/$offset/g" -e "s/%USER%/$user/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/list_files_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
