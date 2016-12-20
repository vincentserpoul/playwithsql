#!/bin/sh
docker rm -f mssqldb;
docker-compose -f ./islatest/mssql/docker-compose-solo.yml up -d;
sleep 5s;

docker exec -it mssqldb /bin/bash -c 'echo "create database entityone_test;" > createdb.sql && /usr/bin/sqlcmd -U sa -P thank5MsSQLforcingMe -i ./createdb.sql;'
docker exec -it mssqldb /bin/bash -c 'echo "create database playwithsql;" > createdb.sql && /usr/bin/sqlcmd -U sa -P thank5MsSQLforcingMe -i ./createdb.sql;'

# to launch the tests benchmark
# ./mssql/container_launch.sh;go test -db=mssql -bench=.  -test.benchtime=3s;docker rm -f mssqldb;