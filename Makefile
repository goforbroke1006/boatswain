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

build/release:
	GOOS=linux   GOARCH=amd64 go build -o "./.build/release/${SERVICE_NAME}-linux-amd64" .
	GOOS=darwin  GOARCH=amd64 go build -o "./.build/release/${SERVICE_NAME}-darwin-amd64" .
	GOOS=windows GOARCH=amd64 go build -o "./.build/release/${SERVICE_NAME}-windows-amd64.exe" .
