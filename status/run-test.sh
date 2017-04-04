#!/bin/bash

for DB in cockroachdb mssql mysql oracle postgres;
do
    ./infra/databases/swarm/$DB/launch-solo.sh;
    go test ./status/ -db=$DB -bench=./status/islatest/$DB -test.benchtime=3s;
    docker service rm pws_$DB;
done;

go test ./status/ -db=sqlite -bench=./status/islatest/sqlite -test.benchtime=3s;rm -f ./status/test.db;
