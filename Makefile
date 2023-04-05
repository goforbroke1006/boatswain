SERVICE_NAME=boatswain

.PHONY: all
all: prepare build test lint

.PHONY: prepare
prepare:
	go mod tidy
	go mod download
	go generate ./...
	go mod tidy

.PHONY: build
build:
	go build -o "${SERVICE_NAME}" .

.PHONY: test
test:
	go test -short ./...

.PHONY: lint
lint: lint/golang

.PHONY: lint/golang
lint/golang:
	golangci-lint run
