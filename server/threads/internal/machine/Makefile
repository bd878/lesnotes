all : commands.proto
	protoc \
		--proto_path=../../../ \
		--go_out=. \
		--go_opt=paths=import \
		--go_opt=module="github.com/bd878/gallery/server/threads/internal/machine" \
		threads/internal/machine/commands.proto