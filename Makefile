.PHONY: vet fmt test build

vet:
	go vet ./...

fmt:
	go fmt -x ./...

test:
	go test ./...

build-dev: vet fmt test
	CGO_ENABLED=0 go build -ldflags "-X main.version=dev" -a -o ./bin/jxv
