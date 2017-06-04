#!/bin/bash

for DB in cockroachdb mssql mysql postgres;
do
    ./infra/databases/swarm/$DB/launch-solo.sh;
    echo "running tests";
    go test ./status/ -db=$DB -sch=islatest -bench=. -test.benchtime=3s;
    go test ./status/ -db=$DB -sch=lateststatus -bench=. -test.benchtime=3s;
    docker service rm pws_$DB;
done;

go test ./status/ -db=sqlite -sch=islatest -bench=. -test.benchtime=3s;rm -f ./status/test.db;
go test ./status/ -db=sqlite -sch=lateststatus -bench=. -test.benchtime=3s;rm -f ./status/test.db;