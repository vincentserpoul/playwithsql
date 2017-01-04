#!/bin/bash

docker-compose -f ./infra/docker_local/cockroachdb/docker-compose-solo.yml down;
docker-compose -f ./infra/docker_local/cockroachdb/docker-compose-solo.yml up -d;
sleep 5s;
docker exec -it roach1 ./cockroach sql --execute="CREATE DATABASE entityone_test;";
docker exec -it roach1 ./cockroach sql --execute="CREATE DATABASE playwithsql;";

# to launch the tests benchmark
# ./cockroachdb/container_launch.sh;go test -db=cockroachdb -bench=.  -test.benchtime=3s;docker rm -f roach1;