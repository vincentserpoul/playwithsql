#!/bin/sh
docker run -dit --name mariadb -e MYSQL_ROOT_PASSWORD=test -p 3306:3306 mariadb
