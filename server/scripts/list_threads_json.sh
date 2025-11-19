#!/usr/bin/bash

# List user threads via json request

stage=${STAGE:-"stage."}
token=${1?"Usage: list_threads_json.sh token user_id parent_id limit offset"}
user=${2?"Usage: list_threads_json.sh token user_id parent_id limit offset"}
parent=${3?"Usage: list_threads_json.sh token user_id parent_id limit offset"}
limit=${4?"Usage: list_threads_json.sh token user_id parent_id limit offset"}
offset=${5?"Usage: list_threads_json.sh token user_id parent_id limit offset"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"limit\":%LIMIT%,\"offset\":%OFFSET%,\"parent\":%PARENT%,\"user_id\":%USER%}}' | 
	sed -e "s/%TOKEN%/$token/g" -e "s/%PARENT%/$parent/g" -e "s/%LIMIT%/$limit/g" -e "s/%OFFSET%/$offset/g" -e "s/%USER%/$user/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/list_threads_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
