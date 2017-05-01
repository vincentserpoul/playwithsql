#!/bin/bash

DB=$1
HOST=$2
LOOPS=$3
MAXCONNS=$4

docker service rm pws_cmd-$DB

docker service create \
    --name pws_cmd-$DB  \
    --restart-condition none \
    --network pws_default \
    vincentserpoul/playwithsql-cmd-status \
    -db=$DB -host=$HOST -loops=$LOOPS -maxconns=$MAXCONNS 

sleep 2s;

while [ $(docker service ls | grep $DB | awk '{print $4}' | grep "0/" | wc -l) != 1 ] ;do 
    sleep 1s;
done;

docker service logs pws_cmd-$DB | awk '{ print $3."," }' >> ./bench/status/swarm/results.log;
docker service rm pws_cmd-$DB;
docker service rm pws_$DB;