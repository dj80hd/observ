# build stage
FROM golang:1.11-alpine AS build
RUN apk add bash ca-certificates git gcc g++ libc-dev make

RUN go get -u golang.org/x/lint/golint

# golang base image has GOPATH=/go
ADD . /go/src/github.com/dj80hd/observ
WORKDIR /go/src/github.com/dj80hd/observ

COPY go.mod .
COPY go.sum .

RUN GO111MODULE=on make test

FROM alpine:latest

COPY --from=build /go/src/github.com/dj80hd/observ/build/observ /observ

ENTRYPOINT ["/observ"]
