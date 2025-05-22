#!/usr/bin/bash

# Signs user up

stage=${STAGE:-"stage."}
name=${NAME:?"Usage: env NAME= PASSWORD= ./signup.sh"}
password=${PASSWORD:?"Usage: env NAME= PASSWORD= ./signup.sh"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%NAME%/$name/g" \
-e "s/%PASSWORD%/$password/g" ./curl/signup.curl |
curl -v -K -
HERE`
result=`eval "$cmd"`
echo $result
