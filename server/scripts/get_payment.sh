#!/usr/bin/bash

# Get payment

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
id=${1:?"Usage: ./get_payment.sh id"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%ID%/$id/g" ./curl/get_payment.curl |
curl -b $cookie -K -
HERE`
result=`eval "$cmd"`
echo $result
