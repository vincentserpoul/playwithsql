#!/bin/bash

initdb () {
    sleep 60s;
    CONTAINER_NAME=$(docker ps --format '{{.Names}}' | grep pws_cockroachdb-0);
    docker exec -i $CONTAINER_NAME ./cockroach sql --execute="CREATE DATABASE entityone_test;";
    docker exec -i $CONTAINER_NAME ./cockroach sql --execute="CREATE DATABASE playwithsql;";
}

removeService () {
    docker service rm pws_cockroachdb-0;
    docker service rm pws_cockroachdb-1;
    docker service rm pws_cockroachdb-2;
}

runService () {
    removeService;
    docker deploy --compose-file ./infra/databases/swarm/cockroachdb/compose-cluster.yml pws;
    initdb;
}

runService;