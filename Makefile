.PHONY: install
install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: protobuf
protobuf: install
	mkdir -p protobuf
	protoc --go_out=protobuf --go_opt=paths=source_relative \
		--go-grpc_out=protobuf --go-grpc_opt=paths=source_relative \
		--proto_path=api/proto \
		api/proto/*/*.proto

.PHONY: build_website
build_website:
	cd website; pnpm install; pnpm run build

.PHONY: swag
swag:
	swag fmt
	swag init -g cmd/web/main.go --ot json -o api --exclude website,api,bin,build,deployment

.PHONY: build
build: protobuf swag build_website
	go build -o bin ./cmd/...