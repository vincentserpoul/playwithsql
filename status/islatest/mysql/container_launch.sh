#!/bin/sh

initdb() {
    sleep 30s;
    docker exec -it mydb mysql -u root -ptest -e 'CREATE DATABASE entityone_test';
    docker exec -it mydb mysql -u root -ptest -e 'CREATE DATABASE playwithsql';
}

removeContainer() {
    docker rm -f mydb;
}

echo "Choose your flavor:"
select flavor in "mysql 8.0.0" "percona 5.7.15" "mariadb 10.1.19"; do
    case $flavor in
        "mysql 8.0.0" ) removeContainer;docker run -dit --name mydb -e MYSQL_ROOT_PASSWORD=test -p 3306:3306 mysql:8.0.0;initdb;break;;
        "percona 5.7.15" ) removeContainer;docker run -dit --name mydb -e MYSQL_ROOT_PASSWORD=test -p 3306:3306 percona:5.7.15;initdb;break;;
        "mariadb 10.1.19" ) removeContainer;docker run -dit --name mydb -e MYSQL_ROOT_PASSWORD=test -p 3306:3306 mariadb:10.1.19;initdb;break;;
    esac
done

# to launch the tests benchmark
# ./mysql/container_launch.sh;go test -db=mysql -bench=.  -test.benchtime=3s;docker rm -f mydb;
