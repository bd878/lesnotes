#!/usr/bin/bash

# Create invoice via json api

stage=${STAGE:-"stage."}
token=${1?:"Usage: create_invoice_json.sh token id"}
id=${2?:"Usage: create_invoice_json.sh token id"}

json=$(echo -n '
{
  "token": "%TOKEN%",
  "req": {
    "id": "%ID%",
    "total": 1000,
    "cart": {
      "items": [
        {
          "type": "premium",
          "item": {
            "expires_at": "2027-01-01T23:59:59Z",
            "cost": 1000,
            "currency": "xtr"
          }
        }
      ]
    }
  }
}' | tr -d ' ' | tr -d '\n' | sed -e 's/"/\\\"/g' -e "s/%TOKEN%/$token/g" -e "s/%ID%/$id/g" )

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" -e "s/%JSON%/$(echo -n $json)/g" ./curl/create_invoice_json.curl | curl -K -
HERE`
# echo $cmd
result=`eval "$cmd"`
echo $result
