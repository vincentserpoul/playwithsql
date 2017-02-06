#!/bin/bash

initdb () {
    sleep 60s;
    CONTAINER_NAME=$(docker ps --format '{{.Names}}' | grep pws_mysql);
    docker exec -i $CONTAINER_NAME mysql -u root -ptest -e 'CREATE DATABASE entityone_test';
    docker exec -i $CONTAINER_NAME mysql -u root -ptest -e 'CREATE DATABASE playwithsql';
}

removeService () {
    docker service rm pws_mysql
}

runService () {
    removeService;
    docker deploy --compose-file ./infra/databases/docker_swarm/mysql/compose-solo.yml pws;
    initdb;
}

runService;