#!/bin/bash

echo "[" > ./bench/status/swarm/results.log
./infra/databases/swarm/cockroachdb/launch-solo.sh && ./bench/status/swarm/run-docker.sh cockroachdb pws_cockroachdb 100 5ms 100;
./infra/databases/swarm/mssql/launch-solo.sh && ./bench/status/swarm/run-docker.sh mssql pws_mssql 100 5ms 3;
./infra/databases/swarm/mysql/launch-solo.sh && ./bench/status/swarm/run-docker.sh mysql pws_mysql 100 5ms 100;
./infra/databases/swarm/mariadb/launch-solo.sh && ./bench/status/swarm/run-docker.sh mariadb pws_mariadb 100 5ms 100;  
./infra/databases/swarm/percona/launch-solo.sh && ./bench/status/swarm/run-docker.sh percona pws_percona 100 5ms 100;  
./infra/databases/swarm/oracle/launch-solo.sh && ./bench/status/swarm/run-docker.sh oracle pws_oracle 100 5ms 100;  
./infra/databases/swarm/postgres/launch-solo.sh && ./bench/status/swarm/run-docker.sh postgres pws_postgres 100 5ms 100;
./bench/status/swarm/run-docker.sh sqlite pws_sqlite 100 10ms 100;

#  remove last comma 
sed -i '$s/,$//' ./bench/status/swarm/results.log;
echo "]" >> ./bench/status/swarm/results.log;

# Generate graph



# Cluster
 # ./infra/databases/swarm/cockroachdb/launch-cluster.sh && ./bench/status/swarm/run-docker.sh cockroachdb pws_cockroachdb-0 2000 && docker service logs -f pws-cmd