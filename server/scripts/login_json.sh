#!/usr/bin/bash

# Logs in user via json api

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
login=${LOGIN:?"Usage: env LOGIN= PASSWORD= ./login.sh"}
password=${PASSWORD:?"Usage: env LOGIN= PASSWORD= ./login.sh"}

json=$(echo -n '{\"login\":\"%LOGIN%\",\"password\":\"%PASSWORD\"}' | sed -e "s/%LOGIN%/$login/g" -e "s/%PASSWORD%/$password/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/login_json.curl |
curl -c $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
