#!/bin/bash
#
# Generates gRPC code given proto definitions
#

if [ $# -ne 0 ]; then
  echo "Usage: $0"
  exit 1;
fi

protoc -I=api --go_out=. --go-grpc_out=. ./api/users.proto

exit 0;

