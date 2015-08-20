#!/bin/bash

export GOROOT="$HOME/cache/gonative/go"
export PATH="$HOME/cache/gonative/go/bin:$PATH"

go run install.go -cmd="gox $ -osarch='windows/amd64 linux/amd64 darwin/amd64'"
