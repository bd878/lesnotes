#!/usr/bin/env bash

# Check status of backend services

declare -A STAGE
declare -A PROD

STAGE["messages"]="https://stage.lesnotes.space/messages/v1/status"
STAGE["files"]="https://stage.lesnotes.space/files/v1/status"
STAGE["users"]="https://stage.lesnotes.space/users/v1/status"
STAGE["search"]="https://stage.lesnotes.space/search/v1/status"
STAGE["billing"]="https://stage.lesnotes.space/billing/v1/status"

PROD["messages"]="https://lesnotes.space/messages/v1/status"
PROD["files"]="https://lesnotes.space/files/v1/status"
PROD["users"]="https://lesnotes.space/users/v1/status"
PROD["search"]="https://lesnotes.space/search/v1/status"
PROD["billing"]="https://lesnotes.space/billing/v1/status"

printf "STAGE:\n"
for i in "${!STAGE[@]}"
do
	res=`curl -s "${STAGE[$i]}"`
	printf "%s: %s\n" $i $res
done

echo

printf "PROD:\n"
for i in "${!PROD[@]}"
do
	res=`curl -s "${PROD[$i]}"`
	printf "%s: %s\n" $i $res
done
