#!/usr/bin/bash

# Cancel payment via json api

stage=${STAGE:-"stage."}
token=${1?:"Usage: cancel_payment_json.sh token id"}
id=${2?:"Usage: cancel_payment_json.sh token id"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"ID\":%ID%}}' | sed -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
sed -e "s/%JSON%/$json/g" ./curl/cancel_payment_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
