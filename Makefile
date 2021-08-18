NAME := sagerconnect
PACKAGE_NAME := github.com/sagernet/sagerconnect

PLATFORM := linux
BUILD_DIR := build
GOBUILD = env CGO_ENABLED=1 $(GO_DIR)go build -v -trimpath -ldflags="-s -w" -o $(BUILD_DIR)
GOBUILD_LINUX = env CGO_ENABLED=1 $(GO_DIR)go build -linkshared -v -trimpath -ldflags="-s -w" -o $(BUILD_DIR)

.PHONY: sagerconnect release
normal: sagerconnect

clean:
	rm -rf $(BUILD_DIR)
	rm -f *.zip
	rm -f *.dat

test:
	# Disable Bloomfilter when testing
	SHADOWSOCKS_SF_CAPACITY="-1" $(GO_DIR)go test -v ./...

sagerconnect:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD)

%.zip: %
	@zip -du $(NAME)-$@ -j $(BUILD_DIR)/$</*
	@zip -du $(NAME)-$@ example/*
	@-zip -du $(NAME)-$@ *.dat
	@echo "<<< ---- $(NAME)-$@"

release: darwin-amd64.zip darwin-arm64.zip linux-386.zip linux-amd64.zip \
	linux-arm.zip linux-armv5.zip linux-armv6.zip linux-armv7.zip linux-armv8.zip \
	linux-mips-softfloat.zip linux-mips-hardfloat.zip linux-mipsle-softfloat.zip linux-mipsle-hardfloat.zip \
	linux-mips64.zip linux-mips64le.zip freebsd-386.zip freebsd-amd64.zip \
	windows-386.zip windows-amd64.zip windows-arm.zip windows-armv6.zip windows-armv7.zip

darwin-amd64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=darwin $(GOBUILD)/$@

darwin-arm64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm64 GOOS=darwin $(GOBUILD)/$@

linux-386:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=386 GOOS=linux $(GOBUILD_LINUX)/$@

linux-amd64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=linux $(GOBUILD_LINUX)/$@

linux-arm:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm GOOS=linux $(GOBUILD_LINUX)/$@

linux-armv5:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm GOOS=linux GOARM=5 $(GOBUILD_LINUX)/$@

linux-armv6:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm GOOS=linux GOARM=6 $(GOBUILD_LINUX)/$@

linux-armv7:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm GOOS=linux GOARM=7 $(GOBUILD_LINUX)/$@

linux-armv8:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm64 GOOS=linux $(GOBUILD_LINUX)/$@

linux-mips-softfloat:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=mips GOMIPS=softfloat GOOS=linux $(GOBUILD_LINUX)/$@

linux-mips-hardfloat:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=mips GOMIPS=hardfloat GOOS=linux $(GOBUILD_LINUX)/$@

linux-mipsle-softfloat:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=mipsle GOMIPS=softfloat GOOS=linux $(GOBUILD_LINUX)/$@

linux-mipsle-hardfloat:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=mipsle GOMIPS=hardfloat GOOS=linux $(GOBUILD_LINUX)/$@

linux-mips64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=mips64 GOOS=linux $(GOBUILD_LINUX)/$@

linux-mips64le:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=mips64le GOOS=linux $(GOBUILD_LINUX)/$@

freebsd-386:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=386 GOOS=freebsd $(GOBUILD_LINUX)/$@

freebsd-amd64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=freebsd $(GOBUILD_LINUX)/$@

windows-386:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=386 GOOS=windows $(GOBUILD)/$@

windows-amd64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=windows $(GOBUILD)/$@

windows-arm:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm GOOS=windows $(GOBUILD)/$@

windows-armv6:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm GOOS=windows GOARM=6 $(GOBUILD)/$@

windows-armv7:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=arm GOOS=windows GOARM=7 $(GOBUILD)/$@
