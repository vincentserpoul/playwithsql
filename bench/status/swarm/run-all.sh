#!/bin/bash

echo "[" > ./bench/status/swarm/results.log
./infra/databases/swarm/cockroachdb/launch-solo.sh && ./bench/status/swarm/run-bench.sh cockroachdb pws_cockroachdb 100 9;
./infra/databases/swarm/mssql/launch-solo.sh && ./bench/status/swarm/run-bench.sh mssql pws_mssql 100 32767;
./infra/databases/swarm/mysql/launch-solo.sh && ./bench/status/swarm/run-bench.sh mysql pws_mysql 100 1000;
./infra/databases/swarm/mariadb/launch-solo.sh && ./bench/status/swarm/run-bench.sh mariadb pws_mariadb 100 1000;  
./infra/databases/swarm/percona/launch-solo.sh && ./bench/status/swarm/run-bench.sh percona pws_percona 100 1000;  
./infra/databases/swarm/oracle/launch-solo.sh && ./bench/status/swarm/run-bench.sh oracle pws_oracle 100 1;  
./infra/databases/swarm/postgres/launch-solo.sh && ./bench/status/swarm/run-bench.sh postgres pws_postgres 100 100;
./bench/status/swarm/run-bench.sh sqlite pws_sqlite 100 100;

#  remove last comma 
sed -i '$s/,$//' ./bench/status/swarm/results.log;
echo "]" >> ./bench/status/swarm/results.log;

# Cluster
 # ./infra/databases/swarm/cockroachdb/launch-cluster.sh && ./bench/status/swarm/run-bench.sh cockroachdb pws_cockroachdb-0 2000 && docker service logs -f pws-cmd