#!/bin/bash

echo "[" > results.log
./infra/databases/docker_swarm/cockroachdb/launch-solo.sh && ./run-docker.sh cockroachdb pws_cockroachdb 100 100;
./infra/databases/docker_swarm/mssql/launch-solo.sh && ./run-docker.sh mssql pws_mssql 100 3;
./infra/databases/docker_swarm/mysql/launch-solo.sh && ./run-docker.sh mysql pws_mysql 100 100;
./infra/databases/docker_swarm/mariadb/launch-solo.sh && ./run-docker.sh mariadb pws_mariadb 100 100;  
./infra/databases/docker_swarm/percona/launch-solo.sh && ./run-docker.sh percona pws_percona 100 100;  
./infra/databases/docker_swarm/oracle/launch-solo.sh && ./run-docker.sh oracle pws_oracle 100 50;  
./infra/databases/docker_swarm/postgres/launch-solo.sh && ./run-docker.sh postgres pws_postgres 100 100;  
# ./run-docker.sh sqlite pws_sqlite 100 1;

#  remove last comma 
sed -i '$s/,$//' results.log;
echo "]" >> results.log;

 # ./infra/databases/docker_swarm/cockroachdb/launch-cluster.sh && ./run-docker.sh cockroachdb pws_cockroachdb-0 2000 && docker service logs -f pws-cmd