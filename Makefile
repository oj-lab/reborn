.PHONY: install
install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	pnpm install -g @openapitools/openapi-generator-cli

.PHONY: protobuf
protobuf:
	mkdir -p protobuf
	protoc --go_out=protobuf --go_opt=paths=source_relative \
		--go-grpc_out=protobuf --go-grpc_opt=paths=source_relative \
		--proto_path=api/proto \
		api/proto/*/*.proto

.PHONY: swag
swag:
	swag fmt
	swag init -g cmd/web/routers/api_v1.go --ot json -o api --instanceName api_v1 \
		--exclude website,api,bin,build,deployment
	openapi-generator-cli generate -i /local/api/api_v1_swagger.json -g typescript-axios -o /local/website/src/api
	rm -rf book/src/api/docs
	mv website/src/api/docs book/src/api

.PHONY: build
build: protobuf swag
	mkdir -p bin
	go build -o bin ./cmd/...

.PHONY: serve_book
serve_book:
	cd book; mdbook serve --open

.PHONY: website
website:
	cd website; pnpm install; pnpm run build