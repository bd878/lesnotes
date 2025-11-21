#!/usr/bin/bash

# Reorder threads via json

stage=${STAGE:-"stage."}
token=${1?:"Usage: reorder_thread_json.sh token id parent next prev"}
id=${2?:"Usage: reorder_thread_json.sh token id parent next prev"}
parent=${3?:"Usage: reorder_thread_json.sh token id parent next prev"}
next=${4?:"Usage: reorder_thread_json.sh token id parent next prev"}
prev=${5?:"Usage: reorder_thread_json.sh token id parent next prev"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"id\":%ID%,\"parent\":%PARENT%,\"next\":%NEXT%,\"prev\":%PREV%}}' |
	sed -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g" -e "s/%PARENT%/$parent/g" -e "s/%NEXT%/$next/g" -e "s/%PREV%/$prev/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
sed -e "s/%JSON%/$json/g" ./curl/reorder_thread_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
