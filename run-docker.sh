#!/bin/bash

DB=$1
HOST=$2
LOOPS=$3

docker service rm pws-cmd-$DB

docker service create \
    --name pws-cmd-$DB  \
    --restart-condition none \
    --network pws_default \
    vincentserpoul/playwithsql-cmd \
    -db=$DB -host=$HOST -loops=$LOOPS

(docker service logs -f pws-cmd-$DB > results-$DB.csv &) 