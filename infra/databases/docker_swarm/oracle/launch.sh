#!/bin/bash

initdb () {
    sleep 60s;
    CONTAINER_NAME=$(docker ps --format '{{.Names}}' | grep pws_oracle);
    docker exec -i $CONTAINER_NAME /bin/bash -c 'ORACLE_HOME="/u01/app/oracle/product/11.2.0/xe" ORACLE_SID="XE" u01/app/oracle/product/11.2.0/xe/bin/sqlplus -s /nolog <<EOF
connect system/oracle
    create user playwithsql identified by "dev";
    grant all privileges to playwithsql;
    create user entityone_test identified by "dev";
    grant all privileges to entityone_test;
quit
EOF';
}

removeContainer () {
    docker service rm pws_oracle
}

runContainer () {
    removeContainer;
    docker deploy --compose-file ./infra/databases/docker_swarm/oracle/compose-solo.yml pws;
    initdb;
}

runContainer