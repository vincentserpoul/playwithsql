#!/bin/sh
docker-compose -f ./infra/docker_local/postgres/docker-compose-solo.yml down;
docker-compose -f ./infra/docker_local/postgres/docker-compose-solo.yml up -d;
sleep 5s;
docker exec -it postgresdb psql -c 'CREATE DATABASE entityone_test;';
docker exec -it postgresdb psql -c 'CREATE DATABASE playwithsql;';

# to launch the tests benchmark
# ./postgres/container_launch.sh;go test -db=postgres -bench=.  -test.benchtime=3s;docker rm -f postgresdb;