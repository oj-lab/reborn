.PHONY: install
install:
	go install github.com/swaggo/swag/cmd/swag@latest
	pnpm install -g @openapitools/openapi-generator-cli

.PHONY: swag
swag:
	swag fmt
	swag init -g internal/routers/api_v1.go --ot json -o api --instanceName api_v1 \
		--exclude website,api,bin,build,deployment --parseDependency -q
	openapi-generator-cli generate -i /local/api/api_v1_swagger.json -g typescript-axios -o /local/website/src/api

.PHONY: build
build: swag
	mkdir -p bin
	go build -o bin/web cmd/main.go

.PHONY: website
website:
	cd website; pnpm install; pnpm run build