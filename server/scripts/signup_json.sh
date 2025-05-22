#!/usr/bin/bash

# Signs user up via json api

name=${NAME:?"Usage: env NAME= PASSWORD= ./signup.sh"}
password=${PASSWORD:?"Usage: env NAME= PASSWORD= ./signup.sh"}

json="{\"name\":\"$name\",\"password\":\"$password\"}"

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%JSON%/$json/g" ./curl/signup_json.curl |
curl -v -K -
HERE`
result=`eval "$cmd"`
echo $result
