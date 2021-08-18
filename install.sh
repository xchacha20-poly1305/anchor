#!/bin/bash

source .github/env.sh

go build -o build/sc -v -linkshared -trimpath -ldflags "-s -w -buildid=" .
sudo cp -f build/sc /usr/local/bin/sc
