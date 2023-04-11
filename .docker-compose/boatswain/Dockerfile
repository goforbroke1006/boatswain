FROM golang:1.20-bullseye

RUN go install github.com/githubnemo/CompileDaemon@v1.4.0

WORKDIR /code/

ADD go.mod .
ADD go.sum .
RUN go mod download -x

COPY . .
RUN go generate ./...

SHELL ["/bin/bash", "-c"]
RUN echo $'#!/bin/bash \n\
\n\
go mod download \n\
go generate ./... \n\
find ./ -name "*.gen.go" -exec chmod 0777 {} \; \n\
find ./ -name "*.mock.go" -exec chmod 0777 {} \; \n\
\n\
appArgs="$*" \n\
\n\
CompileDaemon --build="go build -o /tmp/application -buildvcs=false ." --command="/tmp/application $appArgs" \n\
' > /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
