#!/usr/bin/bash

# List user public messages

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
limit=${LIMIT:-"20"}
offset=${OFFSET:-"0"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%LIMIT%/$limit/g" \
-e "s/%OFFSET%/$offset/g" ./curl/list_public_messages.curl |
curl -b $cookie -s -K -
HERE`
result=`eval "$cmd"`
echo $result
