#!/usr/bin/bash

# Edit existing message

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
update=${UPDATE:-""}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
sed -e "s/%UPDATE%/$update/g" ./curl/edit_message.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
