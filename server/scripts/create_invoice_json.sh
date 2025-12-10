#!/usr/bin/bash

# Create invoice via json api

stage=${STAGE:-"stage."}
token=${1?:"Usage: create_invoice_json.sh token id currency total"}
id=${2?:"Usage: create_invoice_json.sh token id currency total"}
currency=${3?:"Usage: create_invoice_json.sh token id currency total"}
total=${4?:"Usage: create_invoice_json.sh token id currency total"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"id\":\"%ID%\",\"currency\":\"%CURRENCY%\",\"total\":%TOTAL%}}' |
	sed -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g" -e "s/%CURRENCY%/$currency/g" -e "s/%TOTAL%/$total/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/create_invoice_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
