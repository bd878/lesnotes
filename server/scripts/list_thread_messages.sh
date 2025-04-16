#!/usr/bin/bash

# List user thread messages

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
thread_id=${1?"Usage: list_thread_message.sh thread_id limit offset public=-1"}
limit=${2?"Usage: list_thread_messages.sh thread_id limit offset public=-1"}
offset=${3?"Usage: list_thread_messages.sh thread_id limit offset public=-1"}
public=${4:-"-1"}

public_filter=""

if [[ "$public" == "1" ]]; then
	public_filter="public=1"
elif [[ "$public" == "0" ]]; then
	public_filter="public=0"
fi

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%LIMIT%/$limit/g" \
-e "s/%THREAD_ID%/$thread_id/g" \
-e "s/%PUBLIC%/$public_filter/g" \
-e "s/%OFFSET%/$offset/g" ./curl/list_thread_messages.curl |
curl -b $cookie -s -K -
HERE`
result=`eval "$cmd"`
echo $result
