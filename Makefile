SERVICE_NAME=boatswain
IMAGE_NAME=docker.io/goforbroke1006/$(SERVICE_NAME)
IMAGE_TAG=latest

GOLANGCI_LINT_VERSION=v1.52.2

all: prepare build test lint ## Recommended step to prepare project
.PHONY: all

prepare: ## Install dependencies, generate boilerplate code and update go.mod go.sum files
	go mod tidy
	go mod download -x
	go generate ./...
	go mod tidy
.PHONY: prepare

build: ## Compile source code to binary file
	go build -o "$(SERVICE_NAME)" .
.PHONY: build

test: ## Run tests with code coverage print
	@go test -short -coverprofile coverage.out.tmp ./...
	@cat coverage.out.tmp | grep -v ".gen.go" > coverage.out
	@go tool cover -func coverage.out
.PHONY: test

lint: ## Check source code with linter
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@golangci-lint --version
	@golangci-lint run -v .
.PHONY: lint

coverage: ## Run code coverage visual tool to inspect uncovered parts of project
	go test -short -race -coverprofile coverage.out.tmp ./...
	cat coverage.out.tmp | grep -v ".gen.go" > coverage.out
	go tool cover -html ./coverage.out
.PHONY: coverage

benchmark: ## Run benchmark tests and compare with previous results
	go install golang.org/x/perf/cmd/benchstat@latest
	# -run=^#        - skips unit tests
	# -benchtime=10x - adjusts minimum time for each test
	# -benchmem      - print memory usage
	# -cpu=1,2,4     - verify on similar to production settings
	go test -gcflags=-N -bench=. -run=^# -benchtime=10x -benchmem -cpu=1,2,4 ./... | tee .benchmark/new.txt
	benchstat .benchmark/old.txt .benchmark/new.txt
.PHONY: benchmark

image: ## Build image snapshot (latest tag)
	docker build --pull -f ./Dockerfile -t $(IMAGE_NAME):$(IMAGE_TAG) .
	docker push $(IMAGE_NAME):$(IMAGE_TAG)

dev: ## Build local environment docker images
	bash .docker-compose/build-all.sh
.PHONY: dev

start: ## Run local environment
	docker compose down --volumes
	docker compose up -d minimal
.PHONY: start


build/release: ## Build release binaries for Win, MacOS and Linux
	GOOS=linux   GOARCH=amd64 go build -o "./.build/release/$(SERVICE_NAME)-linux-amd64" .
	GOOS=darwin  GOARCH=amd64 go build -o "./.build/release/$(SERVICE_NAME)-darwin-amd64" .
	GOOS=windows GOARCH=amd64 go build -o "./.build/release/$(SERVICE_NAME)-windows-amd64.exe" .


BLUE   = \033[36m
NC 	   = \033[0m

help: ## Prints this help and exits
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "${BLUE}%-30s${NC} %s\n", $$1, $$2}'
.PHONY: help

