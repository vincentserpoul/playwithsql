# playwithsql
implementing as immutable as possible data modelization and benchmarking it on different platforms.

# TODO LIST

- [x] MySQL
- [x] MariaDB
- [x] PerconaDB
- [x] Postgres
- [x] CockroachDB
- [x] SQLite

# The setup

* 1 coreos
* 1 ubuntu 16.04

on ubuntu:
* install golang 1.7.3
* install gcc and sqlite3

```
wget https://storage.googleapis.com/golang/go1.7.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.7.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME
apt-get install -y sqlite3 gcc
go get github.com/vincentserpoul/playwithsql
cd src/github.com/vincentserpoul/playwithsql/
go get ./...
```