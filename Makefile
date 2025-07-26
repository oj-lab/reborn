.PHONY: install swag build build-dev website fmt lint

install:
	go install github.com/swaggo/swag/cmd/swag@latest
	pnpm install -g @openapitools/openapi-generator-cli

swag:
	swag fmt
	swag init -g internal/routers/api_v1.go --ot json -o api --instanceName api_v1 \
		--exclude website,api,bin,build,deployment --parseDependency -q
	openapi-generator-cli generate -i /local/api/api_v1_swagger.json -g typescript-axios -o /local/website/src/api

build:
	mkdir -p bin
	go build -o bin/web cmd/main.go

build-dev: swag
	mkdir -p bin
	go build -o bin/web cmd/main.go

website:
	cd website; pnpm install; pnpm run build

fmt:
	@golines -w .
	@gofumpt -w .

lint:
	@golangci-lint run --fix