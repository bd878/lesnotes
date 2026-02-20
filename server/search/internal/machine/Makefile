all : commands.proto
	protoc \
		--proto_path=../../../ \
		--go_out=. \
		--go_opt=paths=import \
		--go_opt=module="github.com/bd878/gallery/server/search/internal/machine" \
		search/internal/machine/commands.proto