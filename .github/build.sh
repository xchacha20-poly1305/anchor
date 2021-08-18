#!/bin/bash

source .github/env.sh

xgo -v -go go-1.16.7 \
  -dest build \
  -targets "$@" \
  -out sc-$VERSION . || exit 1

chmod -R 777 build
