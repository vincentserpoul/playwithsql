#!/bin/sh
docker rm -f roach1;
docker network create -d bridge roachnet;
docker run -dit --name roach1 --net=roachnet --hostname=roach1 -p 26257:26257 -p 8080:8080 cockroachdb/cockroach  start --insecure;

docker run -dit --name roach2 --net=roachnet --hostname=roach2 cockroachdb/cockroach  start --insecure --join=roach1;
docker run -dit --name roach3 --net=roachnet --hostname=roach3 cockroachdb/cockroach  start --insecure --join=roach1;
docker run -dit --name roach4 --net=roachnet --hostname=roach4 cockroachdb/cockroach  start --insecure --join=roach1;
docker run -dit --name roach5 --net=roachnet --hostname=roach5 cockroachdb/cockroach  start --insecure --join=roach1;
sleep 5s;
docker exec -it roach1 ./cockroach sql --execute="CREATE DATABASE entityone_test;";