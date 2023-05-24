FROM golang:1.20-bullseye

ARG PROJECT_NAME=''

RUN go install github.com/githubnemo/CompileDaemon@v1.4.0

WORKDIR /code/

ADD go.mod .
ADD go.sum .
RUN go mod download -x

# go: downloading golang.org/x/mod v0.4.2
# go: downloading golang.org/x/tools v0.1.1
# go: downloading golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
# go: downloading golang.org/x/sys v0.0.0-20210510120138-977fb7262007

COPY . .
RUN go generate ./...

RUN go run . -- help

RUN printf '#!/bin/bash \n\
\n\
go mod download -x \n\
go generate ./...  \n\
\n\
find ./ -name "*.gen.go" -exec chmod 0777 {} \;  \n\
find ./ -name "*.mock.go" -exec chmod 0777 {} \; \n\
\n\
appArgs="$*" \n\
CompileDaemon --build="go build -o /tmp/application -buildvcs=false ." --command="/tmp/application $appArgs" \n\
' >> /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
