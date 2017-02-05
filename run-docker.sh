#!/bin/bash

DB=$1
HOST="pws_"$1
LOOPS=$2

docker service rm pws-cmd

docker service create \
    --name pws-cmd  \
    --restart-condition none \
    --network pws_default \
    vincentserpoul/playwithsql-cmd \
    -db=$DB -host=$HOST -loops=$LOOPS

# ./run-docker.sh cockroachdb 10000
# ./run-docker.sh mssql 10000
# ./run-docker.sh mysql 10000
# ./run-docker.sh oracle 10000
# ./run-docker.sh postgres 10000
# ./run-docker.sh sqlite 10000