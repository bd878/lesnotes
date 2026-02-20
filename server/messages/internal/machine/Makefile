all : commands.proto
	protoc \
		--proto_path=../../../ \
		--go_out=. \
		--go_opt=paths=import \
		--go_opt=module="github.com/bd878/gallery/server/messages/internal/machine" \
		messages/internal/machine/commands.proto