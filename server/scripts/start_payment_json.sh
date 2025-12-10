#!/usr/bin/bash

# Start payment via json api

stage=${STAGE:-"stage."}
token=${1?:"Usage: start_payment_json.sh token invoice_id id currency total"}
invoice_id=${2?:"Usage: start_payment_json.sh token invoice_id id currency total"}
id=${3?:"Usage: start_payment_json.sh token invoice_id id currency total"}
currency=${4?:"Usage: start_payment_json.sh token invoice_id id currency total"}
total=${5?:"Usage: start_payment_json.sh token invoice_id id currency total"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"id\":%ID%,\"invoice_id\":\"%INVOICE_ID%\",\"currency\":\"%CURRENCY%\",\"total\":%TOTAL%}}' |
	sed -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g" -e "s/%INVOICE_ID%/$invoice_id/g" -e "s/%CURRENCY%/$currency/g" -e "s/%TOTAL%/$total/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/start_payment_json.curl |
curl -K -
HERE`
result=`eval "$cmd"`
echo $result
