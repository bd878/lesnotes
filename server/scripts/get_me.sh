#!/usr/bin/bash

# Get me 

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" ./curl/get_me.curl |
curl -b $cookie -K -
HERE`
result=`eval "$cmd"`
echo $result
