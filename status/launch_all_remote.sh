#!/bin/bash

YOURPRIVATEKEY=$1
COREOSIPBENCH=$2
COREOSIPDB1=$3

# run cockroachdb container
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "
    PATH='/opt/bin:/usr/bin' && cd /home/core/playwithsql &&
    ./infra/databases/docker_local/cockroachdb/container_launch.sh"
ssh -i $YOURPRIVATEKEY core@$COREOSIPBENCH "
    docker rm -f pws-cmd &&
    docker run -t --name pws-cmd vincentserpoul/playwithsql-cmd -db=cockroachdb -host=$COREOSIPDB1 -loops=10000 &&
    docker rm -f pws-cmd"
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "docker rm -f roach1"