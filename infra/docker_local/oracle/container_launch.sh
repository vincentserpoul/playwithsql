#!/bin/bash

docker-compose -f ./infra/docker_local/oracle/docker-compose-solo.yml down;
docker-compose -f ./infra/docker_local/oracle/docker-compose-solo.yml up -d;
sleep 20s;

docker exec -i oracledb /bin/bash -c 'ORACLE_HOME="/u01/app/oracle/product/11.2.0/xe" ORACLE_SID="XE" u01/app/oracle/product/11.2.0/xe/bin/sqlplus -s /nolog <<EOF
connect system/oracle
    create user playwithsql identified by "dev";
    grant all privileges to playwithsql;
    create user entityone_test identified by "dev";
    grant all privileges to entityone_test;
quit
EOF';

# to launch the tests benchmark
# ./oracle/container_launch.sh;go test -db=oracle -bench=.  -test.benchtime=3s;docker rm -f oracledb;