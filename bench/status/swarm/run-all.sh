#!/bin/bash

 export SCH=$1;

# Number of connections is based on n1-standard-1 (1 vCPU, 3.75 GB memory) Machine type on gcloud

echo "[" > ./bench/status/swarm/$SCH/results.log
# until we find out what's wrong with cockroachdb'
./infra/databases/swarm/cockroachdb/launch-solo.sh && ./bench/status/swarm/run-bench.sh cockroachdb cockroachdb $SCH 10000 10;
./infra/databases/swarm/mssql/launch-solo.sh && ./bench/status/swarm/run-bench.sh mssql mssql $SCH 10000 10;
./infra/databases/swarm/mysql/launch-solo.sh && ./bench/status/swarm/run-bench.sh mysql mysql $SCH 10000 10;
./infra/databases/swarm/mariadb/launch-solo.sh && ./bench/status/swarm/run-bench.sh mariadb mariadb $SCH 10000 10;
./infra/databases/swarm/percona/launch-solo.sh && ./bench/status/swarm/run-bench.sh percona percona $SCH 10000 10;
./infra/databases/swarm/postgres/launch-solo.sh && ./bench/status/swarm/run-bench.sh postgres postgres $SCH 10000 10;
./bench/status/swarm/run-bench.sh sqlite pws_sqlite $SCH 100000 1;

#  remove last comma
sed -i '$s/,$//' ./bench/status/swarm/$SCH/results.log;
echo "]" >> ./bench/status/swarm/$SCH/results.log;
