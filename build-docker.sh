#!/bin/bash

set -e -x

GOOS=linux GOARCH=amd64 go build -o playwithsql-cmd ./cmd
docker build ./ --rm -t vincentserpoul/playwithsql-cmd -f ./Dockerfile-cmd
docker push vincentserpoul/playwithsql-cmd