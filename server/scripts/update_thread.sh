#!/usr/bin/bash

# Update a thread.
# Make it public/private

usage="Usage: update_thread.sh id title text name"

stage=${STAGE:-"stage."}
cookie=${COOKIE:-"cookie.txt"}
id=${1?"$usage"}
title=${2?"$usage"}
text=${3?"$usage"}
name=${4?"$name"}

printf "%s %s %s %s\n" $id $title $text $name

cmd=`cat <<HERE
sed -e "s/%STAGE%/$stage/g" \
-e "s/%ID%/$id/g" \
-e "s/%TEXT%/$text/g" \
-e "s/%TITLE%/$title/g" \
-e "s/%NAME%/$name/g" ./curl/update_thread.curl |
curl -b $cookie -K -
HERE`
result=`eval "$cmd"`
echo $result
