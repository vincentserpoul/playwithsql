#!/bin/bash
set -e -x

docker build ./ --rm -t vincentserpoul/playwithsql-cmd-status -f ./Dockerfile-cmd-status
docker push vincentserpoul/playwithsql-cmd-status