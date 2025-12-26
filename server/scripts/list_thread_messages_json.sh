#!/usr/bin/bash

# List user thread messages via json request

stage=${STAGE:-"stage."}
token=${1?"Usage: list_thread_messages_json.sh token user_id thread_id limit offset"}
user=${2?"Usage: list_thread_messages_json.sh token user_id thread_id limit offset"}
thread=${3?"Usage: list_thread_messages_json.sh token user_id thread_id limit offset"}
limit=${4?"Usage: list_thread_messages_json.sh token user_id thread_id limit offset"}
offset=${5?"Usage: list_thread_messages_json.sh token user_id thread_id limit offset"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"user\":%USER%,\"thread\":%THREAD%,\"limit\":%LIMIT%,\"offset\":%OFFSET%}}' |
	sed -e "s/%TOKEN%/$token/g" -e "s/%USER%/$user/g" -e "s/%THREAD%/$thread/g" -e "s/%LIMIT%/$limit/g" -e "s/%OFFSET%/$offset/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/list_thread_messages_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
