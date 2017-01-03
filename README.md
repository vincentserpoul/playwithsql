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
COREOSIP=192.xxx UBUNTUIP=192.xxx YOURPRIVATEKEY=~/.ssh/keykeyxxxx &&
scp -i $YOURPRIVATEKEY ./infra/instantclient-* root@$UBUNTUIP://root/ &&
scp -i $YOURPRIVATEKEY ./infra/install_ubuntu.sh root@$UBUNTUIP://root/ &&
ssh -i $YOURPRIVATEKEY root@$UBUNTUIP "/root/install_ubuntu.sh"
COREOSIP
```

# Launch local status benches

```
./status/launch_all.sh
```