#!/usr/bin/bash

# Send a file

[ -e ".env" ] && source ".env"

stage=${STAGE:-""}
cookie=${COOKIE:-"cookie.txt"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" ./curl/send_file.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
