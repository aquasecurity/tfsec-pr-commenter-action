#!/usr/bin/env bash

export GO111MODULE=auto
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/commenter-linux-amd64 ./cmd/commenter
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/commenter-linux-arm64 ./cmd/commenter
