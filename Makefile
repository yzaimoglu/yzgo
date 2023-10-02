OUT = Backend
VERSION = alpha-0.0.1

.PHONY: run build test

run:
	go run ./cmd/server

build:
	go build -o ./bin/${OUT}_${VERSION} ./cmd/server 

test:
	go test -v ./...