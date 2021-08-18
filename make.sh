#!/bin/bash

if [ ! $(command -v go) ]; then
  if [ -d /usr/lib/go-1.16 ]; then
    export PATH="$PATH:/usr/lib/go-1.16/bin"
  elif [ -d $HOME/.go ]; then
    export PATH="$PATH:$HOME/.go/bin"
  fi
fi

if [ $(command -v go) ]; then
  export PATH="$PATH:$(go env GOPATH)/bin"
fi

make $@