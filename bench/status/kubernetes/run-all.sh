#!/bin/bash

./infra/server/kubernetes/create.sh

echo "[" > ./bench/status/kubernetes/results.log
./infra/databases/kubernetes/cockroachdb/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh cockroachdb cockroachdb 100000 2ms $(($(nproc --al) * 2 + 1));
./infra/databases/kubernetes/mssql/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh mssql mssql 100000 2ms 32767;
./infra/databases/kubernetes/mysql/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh mysql mysql 100000 4ms 1000;
./infra/databases/kubernetes/mariadb/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh mariadb mariadb 100000 4ms 1000;
./infra/databases/kubernetes/percona/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh percona percona 100000 4ms 1000;
./infra/databases/kubernetes/oracle/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh oracle oracle 100000 3ms 1;
./infra/databases/kubernetes/postgres/launch-solo.sh && ./bench/status/kubernetes/run-bench.sh postgres postgres 100000 2ms 100;
./bench/status/kubernetes/run-bench.sh sqlite pws_sqlite 100000 10ms 100;

#  remove last comma 
sed -i '$s/,$//' ./bench/status/kubernetes/results.log;
echo "]" >> ./bench/status/kubernetes/results.log;
