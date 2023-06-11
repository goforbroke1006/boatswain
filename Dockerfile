FROM golang:1.20-bullseye AS builder

WORKDIR /code/

# Warm-up $GOPATH/pkg/mod with libraries that take too much time to install.
RUN go mod init go-mod-warm-up
RUN go get github.com/deepmap/oapi-codegen@v1.13.0
RUN go get github.com/libp2p/go-libp2p@v0.26.3
RUN go get github.com/libp2p/go-libp2p-kad-dht@v0.23.0
RUN go get github.com/libp2p/go-libp2p-pubsub@v0.9.3
RUN go get github.com/multiformats/go-multiaddr@v0.8.0
RUN go get github.com/labstack/echo/v4@v4.10.2
RUN go get github.com/testcontainers/testcontainers-go/modules/compose@v0.19.0

# Copy only deps files and install dependencies.
ADD go.mod .
ADD go.sum .
RUN go mod download -x

# Copy all files and compile to binary file.
COPY . .
RUN go generate ./...
ENV CGO_ENABLED=0
RUN go build -o application -v .



FROM debian:bullseye AS main

RUN apt update
RUN apt install curl -y

WORKDIR /

COPY --from=builder /code/application /application
RUN chmod +x /application

RUN mkdir /db/
COPY ./db/schema.sql /db/schema.sql
RUN ls /db/

ENTRYPOINT [ "/application" ]

HEALTHCHECK --start-period=5s --interval=10s --retries=10 --timeout=5s \
    CMD curl -f http://localhost:8080/readyz || exit 1
