IMAGE="registry.gitlab.com/dj80hd/observ"
SHA:=`git rev-parse --short HEAD`

build: lint
	go build -o build/observ ./cmd/observ

lint:
	go fmt ./...
	golint -set_exit_status $$(go list ./... | grep -v /vendor/)

test: build
	go test -timeout 1s ./...

run: test
	./build/observ

cover: test
	goverage -covermode=set -coverprofile=cov.out `go list ./...`
	gocov convert cov.out | gocov report

coverhtml: cover
	go tool cover -html=cov.out

.PHONY: build build-docker lint run docker publish
