NAME = anchor
TAGS = "with_gvisor"
PARAMS = -v -trimpath -ldflags "-s -w -buildid=" -tags $(TAGS)
MAIN = ./cmd/$(NAME)

.PHONY: build

build:
	go build $(PARAMS) $(MAIN)

fmt:
	@gofumpt -l -w .
	@gofmt -s -w .
	@gci write --custom-order -s standard -s "prefix(github.com/sagernet/)" -s "default" .

fmt_install:
	go install -v mvdan.cc/gofumpt@latest
	go install -v github.com/daixiang0/gci@latest