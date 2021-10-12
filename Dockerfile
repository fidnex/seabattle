FROM golang:1.16-alpine

WORKDIR /go/src/app
COPY . /go/src/app
RUN go build -ldflags="-s -w"  -o '/socket' ./cmd/socket/main.go