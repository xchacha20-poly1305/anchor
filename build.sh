#!/bin/bash

source .github/env.sh

if [[ "$@" == "all" ]]; then

  go run ./infra/build/build.go all

else

  pushd infra/build
  go build -o build .
  popd
  infra/build/build

fi
