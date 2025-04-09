#!/usr/bin/bash

# Authenticates user.
# Checks if given token is valid

[ -e ".env" ] && source ".env"

cookie=${COOKIE:-"cookie.txt"}
stage=${STAGE:-""}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" ./curl/auth.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
