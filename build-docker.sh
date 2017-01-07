#!/bin/bash

set -e -x

GOOS=linux GOARCH=amd64 go build -o playwithsql-cmd ./cmd
