#!/usr/bin/bash

# Deletes user via json api

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
token=${1?:"Usage: delete_user_json.sh token login"}
login=${2?:"Usage: delete_user_json.sh token login"}

json=$(echo -n '{\"token\":\"%TOKEN%\",\"req\":{\"login\":\"%LOGIN%\"}}' | sed -e "s/%TOKEN%/$token/g" -e "s/%LOGIN%/$login/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/delete_user_json.curl |
curl -c $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
