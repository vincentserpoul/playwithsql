#!/bin/bash

set -e -x

GOOS=linux GOARCH=amd64 go build -o playwithsql-cmd-status ./cmd/status/main.go
docker build ./ --rm -t vincentserpoul/playwithsql-cmd-status -f ./Dockerfile-cmd-status
docker push vincentserpoul/playwithsql-cmd-status