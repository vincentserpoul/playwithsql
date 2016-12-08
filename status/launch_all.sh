#!/bin/sh

./islatest/cockroachdb/container_launch.sh;go test -db=cockroachdb -bench=.  -test.benchtime=3s;docker rm -f roach1;
./islatest/mysql/container_launch.sh;go test -db=mysql -bench=.  -test.benchtime=3s;docker rm -f mydb;
./islatest/postgres/container_launch.sh;go test -db=postgres -bench=.  -test.benchtime=3s;docker rm -f postgresdb;
go test -db=sqlite -bench=.  -test.benchtime=3s;rm -f ./test.db;
