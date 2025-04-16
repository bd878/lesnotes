#!/usr/bin/bash

# Sends a message

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
public=${1:-"-1"}

public_filter=""

if [[ "$public" == "1" ]]; then
	public_filter="public=1"
elif [[ "$public" == "0" ]]; then
	public_filter="public=0"
fi

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
sed -e "s/%PUBLIC%/$public_filter/g" ./curl/send_message.curl |
curl -b $cookie -v -K -
HERE`
result=`eval "$cmd"`
echo $result
