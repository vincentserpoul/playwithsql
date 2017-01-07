#!/bin/bash

initdb() {
    sleep 60s;
    docker exec -i mydb mysql -u root -ptest -e 'CREATE DATABASE entityone_test';
    docker exec -i mydb mysql -u root -ptest -e 'CREATE DATABASE playwithsql';
}

removeContainer() {
    docker rm -f mydb;
}

echo "Choose your flavor:"
select flavor in "mysql 8.0.0" "percona 5.7.16" "mariadb 10.1.19"; do
    case $flavor in
        "mysql 8.0.0" ) removeContainer;docker-compose -f ./infra/docker_local/mysql/docker-compose-solo-mysql.yml up -d;initdb;break;;
        "percona 5.7.16" ) removeContainer;docker-compose -f ./infra/docker_local/mysql/docker-compose-solo-percona.yml up -d;initdb;break;;
        "mariadb 10.1.19" ) removeContainer;docker-compose -f ./infra/docker_local/mysql/docker-compose-solo-maria.yml up -d;initdb;break;;
    esac
done

# to launch the tests benchmark
# ./mysql/container_launch.sh;go test -db=mysql -bench=.  -test.benchtime=3s;docker rm -f mydb;
