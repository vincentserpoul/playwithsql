#!/bin/bash

initdb () {
    sleep 60s;
    docker exec -i mydb mysql -u root -ptest -e 'CREATE DATABASE entityone_test';
    docker exec -i mydb mysql -u root -ptest -e 'CREATE DATABASE playwithsql';
}

removeContainer () {
    docker rm -f mydb;
}

runContainer () {
    removeContainer;
    docker-compose -f ./infra/databases/docker_local/mysql/docker-compose-solo-$1.yml up -d;
    initdb;
}

if [ -z "$1" ]
then
    echo "Choose your flavor:"
    select flavor in "mysql 8.0.0" "percona 5.7.16" "mariadb 10.1.20"; do
        case $flavor in
            "mysql 8.0.0" ) runContainer "mysql";break;;
            "percona 5.7.16" ) runContainer "percona";break;;
            "mariadb 10.1.20" ) runContainer "mariadb";break;;
        esac
    done
else
    runContainer $1
fi

# to launch the tests benchmark
# ./mysql/container_launch.sh;go test -db=mysql -bench=.  -test.benchtime=3s;docker rm -f mydb;
