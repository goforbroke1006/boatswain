#!/bin/bash

go mod download
go generate ./...
find ./ -name "*.gen.go" -exec chmod 0777 {} \;
find ./ -name "*.mock.go" -exec chmod 0777 {} \;

appArgs="$*"

CompileDaemon --build="go build -o /tmp/application -buildvcs=false ." --command="/tmp/application $appArgs"
