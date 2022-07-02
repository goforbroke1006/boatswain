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
	go build ./

.PHONY: test
test:
	go test -race ./...

.PHONY: lint
lint: lint/golang

.PHONY: lint/golang
lint/golang:
	golangci-lint run .
