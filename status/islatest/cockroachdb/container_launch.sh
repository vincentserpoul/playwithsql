#!/bin/sh
docker rm -f roach1;
docker run -dit --name roach1  -p 26257:26257 -p 8080:8080 cockroachdb/cockroach:beta-20161110  start --insecure;
sleep 5s;
docker exec -it roach1 ./cockroach sql --execute="CREATE DATABASE entityone_test;";

# to launch the tests benchmark
# go test -db=cockroachdb -bench=.  -test.benchtime=3s
