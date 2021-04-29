#!/usr/bin/env bash

export GO111MODULE=auto

GOOS=linux GOARCH=amd64 go build -o bin/commenter-linux-amd64 ./...
