#!/bin/bash
 
./infra/docker_local/cockroachdb/container_launch.sh;go test ./status/ -db=cockroachdb -bench=./status/islatest/cockroachdb -test.benchtime=3s;docker rm -f roach1;
./infra/docker_local/mysql/container_launch.sh;go test ./status/ -db=mysql -bench=./status/islatest/mysql -test.benchtime=3s;docker rm -f mydb;
./infra/docker_local/postgres/container_launch.sh;go test ./status/ -db=postgres -bench=./status/islatest/postgres -test.benchtime=3s;docker rm -f postgresdb;
./infra/docker_local/mssql/container_launch.sh;go test ./status/ -db=mssql -bench=./status/islatest/mssql -test.benchtime=3s;docker rm -f mssqldb;
./infra/docker_local/oracle/container_launch.sh;go test ./status/ -db=oracle -bench=./status/islatest/oracle -test.benchtime=3s;docker rm -f oracledb;
go test ./status/ -db=sqlite -bench=./status/islatest/sqlite -test.benchtime=3s;rm -f ./status/test.db;
