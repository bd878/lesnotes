all : commands.proto
	protoc \
		--proto_path=../../../../ \
		--go_out=. \
		--go_opt=paths=import \
		--go_opt=module="github.com/bd878/gallery/server/users/internal/controller/distributed" \
		users/internal/controller/distributed/commands.proto