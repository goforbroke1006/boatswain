SERVICE_NAME=boatswain

all: prepare build test lint
.PHONY: all

prepare:
	go mod tidy
	go mod download
	go generate ./...
	go mod tidy
.PHONY: prepare

build:
	go build -o "${SERVICE_NAME}" .
.PHONY: build

test:
	go test -short ./...
.PHONY: test

lint: lint/golang
.PHONY: lint

lint/golang:
	golangci-lint run
.PHONY: lint/golang
