#!/usr/bin/bash

# Signs user up

stage=${STAGE:-"stage."}
login=${LOGIN:?"Usage: env LOGIN= PASSWORD= ./signup.sh"}
password=${PASSWORD:?"Usage: env LOGIN= PASSWORD= ./signup.sh"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%LOGIN%/$login/g" \
-e "s/%PASSWORD%/$password/g" ./curl/signup.curl |
curl -v -K -
HERE`
result=`eval "$cmd"`
echo $result
