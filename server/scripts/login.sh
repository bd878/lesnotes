#!/usr/bin/bash

# Logs in user

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
login=${LOGIN:?"Usage: env LOGIN= PASSWORD= ./login.sh"}
password=${PASSWORD:?"Usage: env LOGIN= PASSWORD= ./login.sh"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%LOGIN%/$login/g" \
-e "s/%PASSWORD%/$password/g" ./curl/login.curl |
curl -c $cookie -K -
HERE`
result=`eval "$cmd"`
echo $result
