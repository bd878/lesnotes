#!/usr/bin/bash

# Get invoice

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
id=${1:?"Usage: ./get_invoice.sh id"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%ID%/$id/g" ./curl/get_invoice.curl |
curl -b $cookie -K -
HERE`
result=`eval "$cmd"`
echo $result
