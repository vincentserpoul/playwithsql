#!/bin/bash

initdb () {
    sleep 60s;
    CONTAINER_NAME=$(docker ps --format '{{.Names}}' | grep pws_postgres);
    docker exec -i $CONTAINER_NAME psql -c 'CREATE DATABASE entityone_test';
    docker exec -i $CONTAINER_NAME psql -c 'CREATE DATABASE playwithsql';
}

removeContainer () {
    docker service rm pws_postgres
}

runContainer () {
    removeContainer;
    docker deploy --compose-file ./infra/databases/docker_swarm/postgres/compose-solo.yml pws;
    initdb;
}

runContainer