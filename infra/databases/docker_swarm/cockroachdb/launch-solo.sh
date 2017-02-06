#!/bin/bash

initdb () {
    sleep 60s;
    CONTAINER_NAME=$(docker ps --format '{{.Names}}' | grep pws_cockroachdb);
    docker exec -i $CONTAINER_NAME ./cockroach sql --execute="CREATE DATABASE entityone_test;";
    docker exec -i $CONTAINER_NAME ./cockroach sql --execute="CREATE DATABASE playwithsql;";
}

removeService () {
    docker service rm pws_cockroachdb
}

runService () {
    removeService;
    docker deploy --compose-file ./infra/databases/docker_swarm/cockroachdb/compose-solo.yml pws;
    initdb;
}

runService;