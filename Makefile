
all:	fmt vet build test

dep:	ensure

ensure:
	dep ensure -v

fmt:
	go fmt ./...

vet:
	go vet -v ./...

build:
	go build ./...

coveralls:
	go test -v -covermode=count -coverprofile=coverage.out ./...
	$GOPATH/bin/goveralls -coverprofile=coverage.out -service=travis-ci

test:
	go test ./...

distclean:
	$(RM) -r vendor