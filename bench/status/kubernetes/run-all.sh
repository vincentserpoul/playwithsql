#!/bin/bash

./infra/server/kubernetes/create.sh

echo "[" > ./bench/status/kubernetes/results.log
./infra/databases/kubernetes/cockroachdb/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh cockroachdb cockroachdb 1000 2ms 100;
./infra/databases/kubernetes/mssql/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh mssql mssql 1000 2ms 100;
./infra/databases/kubernetes/mysql/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh mysql mysql 1000 4ms 1000;
./infra/databases/kubernetes/mariadb/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh mariadb mariadb 1000 4ms 1000;
./infra/databases/kubernetes/percona/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh percona percona 1000 4ms 1000;
./infra/databases/kubernetes/oracle/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh oracle oracle 1000 3ms 100;
./infra/databases/kubernetes/postgres/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh postgres postgres 1000 2ms 100;
./bench/status/kubernetes/run-bench.sh sqlite pws_sqlite 100 10ms 100;

#  remove last comma 
sed -i '$s/,$//' ./bench/status/kubernetes/results.log;
echo "]" >> ./bench/status/kubernetes/results.log;
