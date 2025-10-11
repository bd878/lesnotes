#!/usr/bin/bash

# Search messages by query

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
query=${1?"Usage: search_messages.sh query"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%QUERY%/$query/g" ./curl/search_messages.curl |
curl -b $cookie -s -K -
HERE`
result=`eval "$cmd"`
echo $result
