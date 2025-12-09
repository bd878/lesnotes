#!/usr/bin/bash

# Create invoice via json api

stage=${STAGE:-"stage."}
token=${1?:"Usage: create_invoice_json.sh token language"}
lang=${2?:"Usage: create_invoice_json.sh token language"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"language\":\"%LANG%\"}}' | sed -e "s/%TOKEN%/$token/g" -e "s/%LANG%/$lang/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/create_invoice_json.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
