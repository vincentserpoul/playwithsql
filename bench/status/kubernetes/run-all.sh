#!/bin/bash

export SCH=$1;

# ./infra/server/kubernetes/create.sh
# helm init

# Number of connections is based on n1-standard-1 (1 vCPU, 3.75 GB memory) Machine type on gcloud

echo "[" > ./bench/status/kubernetes/$SCH/results.log
./infra/databases/kubernetes/cockroachdb/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh cockroachdb cockroachdb $SCH 1000 10;
./infra/databases/kubernetes/mssql/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh mssql mssql $SCH 10000 10;
./infra/databases/kubernetes/postgres/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh postgres postgres $SCH 10000 10;
./infra/databases/kubernetes/gcpmysql/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh gcpmysql gcpmysql $SCH 10000 10;
./infra/databases/kubernetes/gcppostgres/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh gcppostgres gcppostgres $SCH 10000 10;
./bench/status/kubernetes/run-bench.sh sqlite pws_sqlite $SCH 10000 1;
./infra/databases/kubernetes/mysql/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh mysql mysql $SCH 10000 10;
./infra/databases/kubernetes/mariadb/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh mariadb mariadb $SCH 10000 10;
./infra/databases/kubernetes/percona/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh percona percona $SCH 10000 10;

#  remove last comma
sed -i '$s/,$//' ./bench/status/kubernetes/$SCH/results.log;
echo "]" >> ./bench/status/kubernetes/$SCH/results.log;
