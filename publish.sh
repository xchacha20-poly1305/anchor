#!/bin/bash

source .github/env.sh

if [ ! -x $(command -v ghr) ]; then
  wget -O ghr.tar.gz https://github.com/tcnksm/ghr/releases/download/v0.14.0/ghr_v0.14.0_linux_amd64.tar.gz
  tar -xvf ghr.tar.gz
  sudo mv ghr*linux_amd64/ghr /usr/local/bin/
  rm -r ghr*linux_amd64
  rm ghr.tar.gz
fi

ghr -delete -t "$GITHUB_TOKEN" -n "v$VERSION" "v$VERSION" build/release
