SERVICE_NAME=boatswain

.PHONY: all
all: dep gen build test lint

.PHONY: dep
dep:
	go mod download

.PHONY: gen
gen:
	go generate ./...

.PHONY: build
build:
	go build -o "${SERVICE_NAME}-demo" ./cmd/demo/
	go build -o "${SERVICE_NAME}-chat" ./cmd/chat/

.PHONY: test
test:
	go test -race ./...

.PHONY: lint
lint: lint/golang

.PHONY: lint/golang
lint/golang:
	golangci-lint run
