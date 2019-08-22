IMAGE="registry.gitlab.com/dj80hd/observ"
SHA:=`git rev-parse --short HEAD`

build: lint
	go build -o build/observ ./cmd/observ

build-docker: lint
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/observ

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

docker: build-docker
	sudo docker build -t $(IMAGE):$(SHA) .

publish: docker
	sudo docker push $(IMAGE):$(SHA)
	sudo docker tag $(IMAGE):$(SHA) $(IMAGE):latest
	sudo docker push $(IMAGE):latest

.PHONY: build build-docker lint run docker publish
