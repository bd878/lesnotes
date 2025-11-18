#!/usr/bin/bash

# Create thread via json api

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
token=${1?:"Usage: create_thread_json.sh token id parent_id next_id prev_id name private"}
id=${2?:"Usage: create_thread_json.sh token id parent_id next_id prev_id name private"}
parent=${3?:"Usage: create_thread_json.sh token id parent_id next_id prev_id name private"}
next=${4?:"Usage: create_thread_json.sh token id parent_id next_id prev_id name private"}
prev=${5?:"Usage: create_thread_json.sh token id parent_id next_id prev_id name private"}
name=${6?:"usage: create_thread_json.sh token id parent_id next_id prev_id name private"}
private=${7?:"usage: create_thread_json.sh token id parent_id next_id prev_id name private"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"id\":%ID%,\"parent\":%PARENT%,\"next\":%NEXT%,\"prev\":%PREV%,\"name\":\"%NAME%\",\"private\":%PRIVATE%}}' |
	sed -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g" -e "s/%PARENT%/$parent/g" -e "s/%NEXT%/$next/g" -e "s/%PREV%/$prev/g" -e "s/%NAME%/$name/g" -e "s/%PRIVATE%/$private/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/create_thread_json.curl |
curl -v -K -
HERE`
result=`eval "$cmd"`
echo $result
