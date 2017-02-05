#!/bin/bash

initdb () {
    sleep 60s;
    CONTAINER_NAME=$(docker ps --format '{{.Names}}' | grep pws_percona);
    docker exec -i $CONTAINER_NAME mysql -u root -ptest -e 'CREATE DATABASE entityone_test';
    docker exec -i $CONTAINER_NAME mysql -u root -ptest -e 'CREATE DATABASE playwithsql';
}

removeContainer () {
    docker service rm pws_percona
}

runContainer () {
    removeContainer;
    docker deploy --compose-file ./infra/databases/docker_swarm/percona/compose-solo.yml pws;
    initdb;
}

runContainer