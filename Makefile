all: generate build

.PHONY: generate
generate:
	go generate ./...
	rm -rf web/src/app/dataaccess/api
	openapi-generator generate -i api/swagger.json -g typescript-fetch -o web/src/app/dataaccess/api

.PHONY: build-judge
build-judge:
	go build -o build/judge cmd/judge/*.go

.PHONY: build-worker
build-worker:
	go build -o build/worker cmd/worker/*.go

.PHONY: build
build: build-judge build-worker

.PHONY: run-judge
run-judge:
	go run cmd/judge/*.go

.PHONY: run-worker
run-worker:
	go run cmd/worker/*.go