all : commands.proto
	protoc \
		--proto_path=../../../ \
		--go_out=. \
		--go_opt=paths=import \
		--go_opt=module="github.com/bd878/gallery/server/users/internal/machine" \
		users/internal/machine/commands.proto