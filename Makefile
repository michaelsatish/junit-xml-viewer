.PHONY: vet fmt build

vet:
	go vet ./...

fmt:
	go fmt ./...

build:
	CGO_ENABLED=0 go build -a -o ./bin/jxv
