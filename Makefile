
all:	fmt vet test server client

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

server:
	go build -o build/coap_server cmd/server/main.go

client:
	go build -o build/coap_client cmd/client/main.go

clean:
	$(RM) -r build

.PHONY: clean fmt vet test server client