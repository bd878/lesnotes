#!/usr/bin/bash

# Update me

usage="Usage: update_me.sh login"

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
login=${1?"$usage"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%LOGIN%/$login/g" ./curl/update_me.curl |
curl -b $cookie -K -
HERE`
result=`eval "$cmd"`
echo $result
