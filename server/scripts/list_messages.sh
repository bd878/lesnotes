#!/usr/bin/bash

# List user messages

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
limit=${1?"Usage: list_messages.sh limit offset public=-1"}
offset=${2?"Usage: list_messages.sh limit offset public=-1"}
public=${3:-"-1"}

public_filter=""

if [[ "$public" == "1" ]]; then
	public_filter="public=1"
elif [[ "$public" == "0" ]]; then
	public_filter="public=0"
fi

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%LIMIT%/$limit/g" \
-e "s/%PUBLIC%/$public_filter/g" \
-e "s/%OFFSET%/$offset/g" ./curl/list_messages.curl |
curl -b $cookie -s -K -
HERE`
result=`eval "$cmd"`
echo $result
