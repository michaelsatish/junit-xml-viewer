.PHONY: vet fmt test build

vet:
	go vet ./...

fmt:
	go fmt -x ./...

test:
	go test ./...

build-dev: vet fmt test
	CGO_ENABLED=0 go build -ldflags "-X main.version=dev" -a -o ./bin/jxv

build: vet fmt test
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=${VERSION}" -a -o ./bin/jxv
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=${VERSION}" -a -o ./bin/jxv.exe
