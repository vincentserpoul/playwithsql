#!/bin/bash

set -e -x

docker run -it -v "${GOPATH}":/gopath -v "$(pwd)":/app -e "GOPATH=/gopath" -w /app vincentserpoul/playwithsql-build sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o playwithsqlcmd ./cmd'

docker build ./ --rm -t playwithsqlcmd -f ./Dockerfile-cmd