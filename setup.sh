#!/bin/bash

set -e

go version

go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2
