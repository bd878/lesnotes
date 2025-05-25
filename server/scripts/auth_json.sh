#!/usr/bin/bash

# Authenticates user via json api.
# Checks if given token is valid

stage=${STAGE:-"stage."}
token=${1?:"Usage: auth_json.sh token"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{}}' | sed -e "s/%TOKEN%/$token/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/auth_json.curl |
curl -v -K - 
HERE`
result=`eval "$cmd"`
echo $result
