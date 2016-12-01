#!/bin/sh
docker rm -f postgresdb;
docker run -dit --name postgresdb -e POSTGRES_USER=root -e POSTGRES_PASSWORD=test -p 5432:5432 postgres:9.6.1-alpine;
sleep 5s;
docker exec -it postgresdb psql -c 'CREATE DATABASE entityone_test;';
docker exec -it postgresdb psql -c 'CREATE DATABASE playwithsql;';

# to launch the tests benchmark
# ./postgres/container_launch.sh;go test -db=postgres -bench=.  -test.benchtime=3s;docker rm -f postgresdb;