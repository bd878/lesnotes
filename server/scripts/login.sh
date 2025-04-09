#!/usr/bin/bash

# Logs in user

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
name=${NAME:?"Usage: env NAME= PASSWORD= ./login.sh"}
password=${PASSWORD:?"Usage: env NAME= PASSWORD= ./login.sh"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%NAME%/$name/g" \
-e "s/%PASSWORD%/$password/g" ./curl/login.curl |
curl -c $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
