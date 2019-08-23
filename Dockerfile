#RUN apk add bash ca-certificates git gcc g++ libc-dev make
FROM golang:1.11-alpine AS build
RUN apk add git

ADD . /go/src/github.com/dj80hd/observ
WORKDIR /go/src/github.com/dj80hd/observ
RUN CGO_ENABLED=0 GO111MODULE=on go build -o build/observ ./cmd/observ


FROM scratch
COPY --from=build /go/src/github.com/dj80hd/observ/build/observ /observ
ENTRYPOINT ["/observ"]
