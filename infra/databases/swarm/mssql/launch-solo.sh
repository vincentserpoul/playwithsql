#!/bin/bash

initdb () {
    sleep 45s;
    CONTAINER_NAME=$(docker ps --format '{{.Names}}' | grep pws_mssql);
    docker exec -i $CONTAINER_NAME /bin/bash -c 'echo "create database entityone_test;" > createdb.sql && /opt/mssql-tools/bin/sqlcmd -U sa -P thank5MsSQLforcingMe -i ./createdb.sql;';
    docker exec -i $CONTAINER_NAME /bin/bash -c 'echo "create database playwithsql;" > createdb.sql && /opt/mssql-tools/bin/sqlcmd -U sa -P thank5MsSQLforcingMe -i ./createdb.sql;';
}

removeService () {
    docker service rm pws_mssql
}

runService () {
    removeService;
    docker deploy --compose-file ./infra/databases/swarm/mssql/compose-solo.yml pws;
    initdb;
}

runService;