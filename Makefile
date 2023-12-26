VERSION := $(shell cat VERSION)
COMMIT_HASH := $(shell git rev-parse HEAD)
PROJECT_NAME := cato

all: generate build-all

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

.PHONY: build-web
build-web:
	cd web && npm run build && cd ..

.PHONY: build-linux-amd64
build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build \
		-ldflags "-X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
		-o build/$(PROJECT_NAME)_linux_amd64 cmd/$(PROJECT_NAME)/*.go

.PHONY: build-linux-arm64
build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build \
		-ldflags "-X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
		-o build/$(PROJECT_NAME)_linux_arm cmd/$(PROJECT_NAME)/*.go

.PHONY: build-macos-amd64
build-macos-amd64:
	GOOS=darwin GOARCH=amd64 go build \
		-ldflags "-X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
		-o build/$(PROJECT_NAME)_macos_amd64 cmd/$(PROJECT_NAME)/*.go

.PHONY: build-macos-arm64
build-macos-arm64:
	GOOS=darwin GOARCH=arm64 go build \
		-ldflags "-X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
		-o build/$(PROJECT_NAME)_macos_arm64 cmd/$(PROJECT_NAME)/*.go

.PHONY: build-windows-amd64
build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build \
		-ldflags "-X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
		-o build/$(PROJECT_NAME)_windows_amd64.exe cmd/$(PROJECT_NAME)/*.go

.PHONY: build-windows-arm64
build-windows-arm64:
	GOOS=windows GOARCH=amd64 go build \
		-ldflags "-X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
		-o build/$(PROJECT_NAME)_windows_arm64.exe cmd/$(PROJECT_NAME)/*.go

.PHONY: build-all
build-all:
	make build-web
	make build-linux-amd64
	make build-linux-arm64
	make build-macos-amd64
	make build-macos-arm64
	make build-windows-amd64
	make build-windows-arm64

.PHONY: build
build:
	go build \
		-ldflags "-X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
		-o build/$(PROJECT_NAME) \
		cmd/$(PROJECT_NAME)/*.go

.PHONY: clean
clean:
	rm -rf build/

.PHONY: run
run:
	go run cmd/$(PROJECT_NAME)/*.go

.PHONY: lint
lint:
	golangci-lint run ./... 