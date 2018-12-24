
all:	fmt vet build test

fmt:
	go fmt ./...

vet:
	go vet ./...

build:
	go build ./...

vendor:
	go mod vendor

coveralls:
	go test -v -covermode=count -coverprofile=coverage.out ./...
	${GOPATH}/bin/goveralls -coverprofile=coverage.out -service=travis-ci

test:
	go test ./...
