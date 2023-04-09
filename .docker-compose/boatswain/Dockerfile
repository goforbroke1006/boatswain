FROM golang:1.20-bullseye AS builder

RUN go install github.com/githubnemo/CompileDaemon@v1.4.0

WORKDIR /code/

ADD go.mod .
ADD go.sum .
RUN go mod download -x

COPY . .
ENV CGO_ENABLED=1
RUN go build -o application .



FROM debian:bullseye AS main

WORKDIR /

COPY --from=builder /code/application /application
RUN chmod +x /application

RUN mkdir /db/
COPY ./db/schema.sql /db/schema.sql
RUN ls /db/

ENTRYPOINT [ "/application" ]
