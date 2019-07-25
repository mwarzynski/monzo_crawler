export CGO_ENABLED=0

all: test build docker

test:
	@go test -cover ./...

build-cli:
	@go build -o bin/crawler-cli cmd/cli/main.go

build-http:
	@go build -o bin/crawler-http-server cmd/http-server/main.go

build: build-cli build-http

docker:
	@docker build -t mwarzynski/crawler .
