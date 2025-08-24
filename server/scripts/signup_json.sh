#!/usr/bin/bash

# Signs user up via json api

stage=${STAGE:-"stage."}
login=${LOGIN:?"Usage: env LOGIN= PASSWORD= ./signup_json.sh"}
password=${PASSWORD:?"Usage: env LOGIN= PASSWORD= ./signup_json.sh"}

json=$(echo -n '{\"login\":\"%LOGIN%\",\"password\":\"%PASSWORD%\"}' | sed -e "s/%LOGIN%/$login/g" -e "s/%PASSWORD%/$password/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/signup_json.curl |
curl -v -K -
HERE`
result=`eval "$cmd"`
echo $result
