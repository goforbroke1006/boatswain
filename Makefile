SERVICE_NAME=boatswain

all: prepare build test lint
.PHONY: all

prepare:
	go mod tidy
	go mod download -x
	go generate ./...
	go mod tidy
.PHONY: prepare

build:
	go build -o "$(SERVICE_NAME)" .
.PHONY: build

test:
	go test -short ./...
.PHONY: test

lint: lint/golang
.PHONY: lint

lint/golang:
	golangci-lint run .
.PHONY: lint/golang

image:
	docker build -f ./Dockerfile -t docker.io/goforbroke1006/$(SERVICE_NAME):latest .
	docker push docker.io/goforbroke1006/$(SERVICE_NAME):latest

dev:
	bash .docker-compose/build-all.sh
.PHONY: dev

build/release:
	GOOS=linux   GOARCH=amd64 go build -o "./.build/release/$(SERVICE_NAME)-linux-amd64" .
	GOOS=darwin  GOARCH=amd64 go build -o "./.build/release/$(SERVICE_NAME)-darwin-amd64" .
	GOOS=windows GOARCH=amd64 go build -o "./.build/release/$(SERVICE_NAME)-windows-amd64.exe" .
