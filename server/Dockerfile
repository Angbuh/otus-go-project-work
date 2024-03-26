FROM golang:1.21-alpine

RUN apk update && apk add musl musl-dev gcc

WORKDIR /program
COPY ./ ./

RUN go build -o prog ./cmd/app/*.go