#!/usr/bin/bash

# Update a message.
# Make it public/private

usage="Usage: update_message.sh message_id text public=(1|0)"

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
id=${1?"$usage"}
text=${2?"$usage"}
public=${3:-"-1"}

printf "%s %s %s\n" $id $text $public

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
sed -e "s/%ID%/$id/g" \
sed -e "s/%TEXT%/$text/g" \
sed -e "s/%PUBLIC%/$public/g" ./curl/update_message.curl |
curl  --trace-ascii /dev/stdout -b $cookie -K -
HERE`
result=`eval "$cmd"`
echo $result
