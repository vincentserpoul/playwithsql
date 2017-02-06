#!/bin/bash

DB=$1
HOST=$2
LOOPS=$3

docker service rm pws-cmd

docker service create \
    --name pws-cmd  \
    --restart-condition none \
    --network pws_default \
    vincentserpoul/playwithsql-cmd \
    -db=$DB -host=$HOST -loops=$LOOPS

# ./infra/databases/docker_swarm/cockroachdb/launch-solo.sh && ./run-docker.sh cockroachdb pws_cockroachdb 2000 && docker service logs -f pws-cmd
# ./infra/databases/docker_swarm/mssql/launch-solo.sh && ./run-docker.sh mssql pws_mssql 2000 && docker service logs -f pws-cmd
# ./infra/databases/docker_swarm/mysql/launch-solo.sh && ./run-docker.sh mysql pws_mysql 2000 && docker service logs -f pws-cmd
# ./infra/databases/docker_swarm/mariadb/launch-solo.sh && ./run-docker.sh mariadb pws_mariadb 2000 && docker service logs -f pws-cmd
# ./infra/databases/docker_swarm/percona/launch-solo.sh && ./run-docker.sh percona pws_percona 2000 && docker service logs -f pws-cmd
# ./infra/databases/docker_swarm/oracle/launch-solo.sh && ./run-docker.sh oracle pws_oracle 2000 && docker service logs -f pws-cmd
# ./infra/databases/docker_swarm/postgres/launch-solo.sh && ./run-docker.sh postgres pws_postgres 2000 && docker service logs -f pws-cmd
# ./run-docker.sh sqlite pws_sqlite 2000 && docker service logs -f pws-cmd

# # ./infra/databases/docker_swarm/cockroachdb/launch-cluster.sh && ./run-docker.sh cockroachdb pws_cockroachdb-0 2000 && docker service logs -f pws-cmd