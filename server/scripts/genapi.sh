#!/bin/bash
#
# Generates gRPC code given proto definitions
#

if [ $# -ne 0 ]; then
	echo "Usage: $0"
	exit 1;
fi

protoc ./protos/*proto \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=import \
	--go-grpc_opt=paths=import \
	--go_opt=module="github.com/bd878/gallery/server" \
	--go-grpc_opt=module="github.com/bd878/gallery/server"

exit 0;

