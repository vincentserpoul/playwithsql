[![Go Report Card](https://goreportcard.com/badge/github.com/vincentserpoul/playwithsql)](https://goreportcard.com/report/github.com/vincentserpoul/playwithsql)

# playwithsql
implementing as immutable as possible data modelization and benchmarking it on different platforms.

## Disclaimer

the benchmark comparison is for very specific use case:
* Golang
* Containerized DBs
* Specific schemas
* Used configurations

Hence, they can't be used to affirm that this or this db is better.
*The context matters!*

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

* 2 coreos machines

# Rebuilding the docker image

Download [Oracle instant client](http://www.oracle.com/technetwork/topics/linuxx86-64soft-092277.html) for Oracle (basic and SDK) to the infra/build folder

```
./build-docker.sh
```

# Launch local status benches

```
./status/launch_all.sh
```

# Launch remote tests

install docker-compose on coreos
```
COREOSIPBENCH=192.xxx COREOSIPDB1=192.xxx YOURPRIVATEKEY=~/.ssh/keykeyxxxx &&
scp -i $YOURPRIVATEKEY ./infra/build/prepare_coreos.sh core@$COREOSIPBENCH://home/core &&
ssh -i $YOURPRIVATEKEY core@$COREOSIPBENCH "chmod +x prepare_coreos.sh" &&
ssh -i $YOURPRIVATEKEY core@$COREOSIPBENCH "./prepare_coreos.sh"
scp -i $YOURPRIVATEKEY ./infra/build/prepare_coreos.sh core@$COREOSIPDB1://home/core &&
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "chmod +x prepare_coreos.sh" &&
ssh -i $YOURPRIVATEKEY core@$COREOSIPDB1 "./prepare_coreos.sh"
```

```
./status/launch_all_remote.sh
```