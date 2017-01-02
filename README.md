[![Go Report Card](https://goreportcard.com/badge/github.com/vincentserpoul/playwithsql)](https://goreportcard.com/report/github.com/vincentserpoul/playwithsql)

# playwithsql
implementing as immutable as possible data modelization and benchmarking it on different platforms.

# TODO LIST

- [x] MySQL
- [x] MariaDB
- [x] PerconaDB
- [x] Postgres
- [x] CockroachDB
- [x] SQLite
- [x] Microsoft SQL Server
- [x] Oracle

# The setup

* 1 coreos
* 1 ubuntu 16.04

# Prerequisites 

Download [Oracle instant client](http://www.oracle.com/technetwork/topics/linuxx86-64soft-092277.html) for Oracle (basic and SDK) to the infra folder

```
scp -i ~/.ssh/YOURPRIVATEKEY ./infra/instantclient-* root@YOURUBUNTUIP://root/
scp -i ~/.ssh/YOURPRIVATEKEY ./infra/install_ubuntu.sh root@YOURUBUNTUIP://root/
ssh -i ~/.ssh/YOURPRIVATEKEY root@YOURUBUNTUIP "/root/install_ubuntu.sh"
```

# Launch local status benches

```
./status/launch_all.sh
```