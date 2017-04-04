#!/bin/bash

initdb () {
    sleep 60s;
    CONTAINER_NAME=$(docker ps --format '{{.Names}}' | grep pws_postgres);
    docker exec -i $CONTAINER_NAME psql -c 'CREATE DATABASE entityone_test';
    docker exec -i $CONTAINER_NAME psql -c 'CREATE DATABASE playwithsql';
}

removeService () {
    docker service rm pws_postgres
}

runService () {
    removeService;
    docker deploy --compose-file ./infra/databases/swarm/postgres/compose-solo.yml pws;
    initdb;
}

runService;