#!/usr/bin/bash

# Deletes a message

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
id=${1?"Usage: delete_message.sh message_id"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
sed -e "s/%ID%/$id/g" ./curl/delete_message.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
