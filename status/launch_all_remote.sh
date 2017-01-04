#!/bin/bash

# run cockroachdb container
ssh -i $YOURPRIVATEKEY core@$COREOSIP "
    cd /home/core/playwithsql &&
    ./infra/docker_local/cockroachdb/container_launch.sh"
ssh -i $YOURPRIVATEKEY root@$UBUNTUIP "
    cd root/src/github.com/vincentserpoul/playwithsql/ && 
    go test ./status/ -db=cockroachdb -host=$COREOSIP -bench=./status/islatest/cockroachdb -test.benchtime=3s &&
    docker rm -f roach1"