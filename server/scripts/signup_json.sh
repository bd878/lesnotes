#!/usr/bin/bash

# Signs user up via json api

stage=${STAGE:-"stage."}
name=${NAME:?"Usage: env NAME= PASSWORD= ./signup.sh"}
password=${PASSWORD:?"Usage: env NAME= PASSWORD= ./signup.sh"}

json=$(echo -n '{\"name\":\"%NAME%\",\"password\":\"%PASSWORD\"}' | sed -e "s/%NAME%/$name/g" -e "s/%PASSWORD%/$password/g")

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/signup_json.curl |
curl -v -K -
HERE`
result=`eval "$cmd"`
echo $result
