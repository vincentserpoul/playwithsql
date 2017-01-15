#!/bin/bash

YOURPRIVATEKEY=$1
LOOPS=$2
COREOSIPBENCH=$3
COREOSIPDB1=$4


# run cockroachdb container
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "
    PATH='/opt/bin:/usr/bin' && cd /home/core/playwithsql &&
    ./infra/databases/docker_local/cockroachdb/container_launch.sh"
ssh -i $YOURPRIVATEKEY core@$COREOSIPBENCH "
    (docker rm -f pws-cmd  || true) &&
    docker run -t --name pws-cmd vincentserpoul/playwithsql-cmd -db=cockroachdb -host=$COREOSIPDB1 -loops=$LOOPS &&
    docker rm -f pws-cmd"
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "docker rm -f roach1"

# run mssql container
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "
    PATH='/opt/bin:/usr/bin' && cd /home/core/playwithsql &&
    ./infra/databases/docker_local/mssql/container_launch.sh"
ssh -i $YOURPRIVATEKEY core@$COREOSIPBENCH "
    (docker rm -f pws-cmd  || true) &&
    docker run -t --name pws-cmd vincentserpoul/playwithsql-cmd -db=mssql -host=$COREOSIPDB1 -loops=$LOOPS &&
    docker rm -f pws-cmd"
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "docker rm -f mssqldb"

# run mysql container
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "
    PATH='/opt/bin:/usr/bin' && cd /home/core/playwithsql &&
    ./infra/databases/docker_local/mysql/container_launch.sh mysql"
ssh -i $YOURPRIVATEKEY core@$COREOSIPBENCH "
    (docker rm -f pws-cmd  || true) &&
    docker run -t --name pws-cmd vincentserpoul/playwithsql-cmd -db=mysql -host=$COREOSIPDB1 -loops=$LOOPS &&
    docker rm -f pws-cmd"
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "docker rm -f mydb"

# run percona container
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "
    PATH='/opt/bin:/usr/bin' && cd /home/core/playwithsql &&
    ./infra/databases/docker_local/mysql/container_launch.sh percona"
ssh -i $YOURPRIVATEKEY core@$COREOSIPBENCH "
    (docker rm -f pws-cmd  || true) &&
    docker run -t --name pws-cmd vincentserpoul/playwithsql-cmd -db=percona -host=$COREOSIPDB1 -loops=$LOOPS &&
    docker rm -f pws-cmd"
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "docker rm -f mydb"

# run maria container
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "
    PATH='/opt/bin:/usr/bin' && cd /home/core/playwithsql &&
    ./infra/databases/docker_local/mysql/container_launch.sh mariadb"
ssh -i $YOURPRIVATEKEY core@$COREOSIPBENCH "
    (docker rm -f pws-cmd  || true) &&
    docker run -t --name pws-cmd vincentserpoul/playwithsql-cmd -db=mariadb -host=$COREOSIPDB1 -loops=$LOOPS &&
    docker rm -f pws-cmd"
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "docker rm -f mydb"

# run oracle container
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "
    PATH='/opt/bin:/usr/bin' && cd /home/core/playwithsql &&
    ./infra/databases/docker_local/oracle/container_launch.sh"
ssh -i $YOURPRIVATEKEY core@$COREOSIPBENCH "
    (docker rm -f pws-cmd  || true) &&
    docker run -t --name pws-cmd vincentserpoul/playwithsql-cmd -db=oracle -host=$COREOSIPDB1 -loops=$LOOPS &&
    docker rm -f pws-cmd"
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "docker rm -f oracledb"

# run postgres container
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "
    PATH='/opt/bin:/usr/bin' && cd /home/core/playwithsql &&
    ./infra/databases/docker_local/postgres/container_launch.sh"
ssh -i $YOURPRIVATEKEY core@$COREOSIPBENCH "
    (docker rm -f pws-cmd  || true) &&
    docker run -t --name pws-cmd vincentserpoul/playwithsql-cmd -db=postgres -host=$COREOSIPDB1 -loops=$LOOPS &&
    docker rm -f pws-cmd"
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "docker rm -f postgresdb"
