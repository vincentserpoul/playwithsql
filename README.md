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

* 1x coreos 16 GB Memory / 160 GB Disk / SGP1 - CoreOS 1185.3.0 (stable)
* 1x ubuntu 16 GB Memory / 160 GB Disk / SGP1 - Ubuntu 16.04.1 x64

on ubuntu:

```
wget https://storage.googleapis.com/golang/go1.7.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.7.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME
apt-get install sqlite3 gcc
go get github.com/vincentserpoul/playwithsql
cd src/github.com/vincentserpoul/playwithsql/
go get ./...
```