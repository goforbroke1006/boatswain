FROM golang:1.20-bullseye AS builder

# install compiled lib to reduce build time
RUN #go install -a github.com/mattn/go-sqlite3@v1.14.16

WORKDIR /code/

ADD go.mod .
ADD go.sum .
RUN go mod download -x

COPY . .
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
