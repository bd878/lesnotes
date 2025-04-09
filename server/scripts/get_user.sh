#!/usr/bin/bash

# Gets user by id

[ -e ".env" ] && source ".env"

stage=${STAGE:-""}
cookie=${COOKIE:-"cookie.txt"}
id=${ID:?"Usage: env ID=<userId> ./get_user.sh"}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%ID%/$id/g" ./curl/get_user.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
