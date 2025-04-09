#!/usr/bin/bash

# Logout user

[ -e ".env" ] && source ".env"

stage=${STAGE:-""}
cookie=${COOKIE:-"cookie.txt"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" ./curl/logout.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
