#!/usr/bin/bash

# Authenticates user.
# Checks if given token is valid

cookie=${COOKIE:-"cookie.txt"}
stage=${STAGE:-"stage."}

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" ./curl/auth.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
