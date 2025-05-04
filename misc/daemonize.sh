#!/usr/bin/bash
# 
# Daemonizes given script.
# Sets umask, changes working dir
# and substitutes current process
# with given script
#

[ $# -ne 1 ] && echo "Usage: daemonize.sh script" && exit 1

umask 0
setsid -f sh -c $1

echo "Process id: `pidof \"$1\"`"

exit 0;

