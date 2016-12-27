#!/bin/sh
docker rm -f oracledb;
docker-compose -f ./islatest/oracle/docker-compose-solo.yml up -d;
sleep 10s;

docker exec -it oracledb /bin/bash -c 'ORACLE_HOME="/u01/app/oracle/product/11.2.0/xe" ORACLE_SID="XE" u01/app/oracle/product/11.2.0/xe/bin/sqlplus -s /nolog <<EOF
connect system/oracle
create user playwithsql identified by "dev";
create user entityone_test identified by "dev";
quit
EOF';

# to launch the tests benchmark
# ./oracle/container_launch.sh;go test -db=oracle -bench=.  -test.benchtime=3s;docker rm -f oracledb;