all: generate build

.PHONY: generate
generate:
	rm -rf internal/handlers/http/rpc/rpcclient
	rm -rf internal/handlers/http/rpc/rpcserver
	rm -rf web/src/app/dataaccess/api

	go get gitlab.com/pjrpc/pjrpc/cmd/genpjrpc@v0.4.0
	go get github.com/google/wire/cmd/wire@v0.5.0
	
	go generate ./...
	openapi-generator generate -i api/swagger.json -g typescript-fetch -o web/src/app/dataaccess/api

	go mod tidy

.PHONY: build
build:
	go build -o build/cato cmd/cato/*.go

.PHONY: run
run:
	go run cmd/cato/*.go
