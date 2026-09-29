NAME = anchor
VERSION = v0.8.0
PARAMS = -v -trimpath -ldflags "-s -w -buildid= -X main.version=$(VERSION)"
MAIN = ./cmd/$(NAME)

.PHONY: build

build:
	CGO_ENABLED=0 go build $(PARAMS) $(MAIN)

fmt:
	@golangci-lint fmt

fmt_install:
	go install -v github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

test:
	go test -v -count=1 ./...