[![Go Report Card](https://goreportcard.com/badge/github.com/vincentserpoul/playwithsql)](https://goreportcard.com/report/github.com/vincentserpoul/playwithsql)

# playwithsql

implementing as immutable as possible data modelization and benchmarking it on different platforms.

## Disclaimer

the benchmark comparison is for very specific use case:
* Golang 1.8
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

- [] swarm cluster - solo db container
- [] kubernetes cluster - solo db container
- [] swarm cluster - cluster db containers
- [] kubernetes cluster - cluster db containers

# Rebuilding the docker image

Download [Oracle instant client](http://www.oracle.com/technetwork/topics/linuxx86-64soft-092277.html) for Oracle (basic and SDK) to the infra/build folder

```
./build-docker-cmd-status.sh
```

# Launch tests

```
./status/run-test.sh
```

# Launch local status benches

```
./bench/status/swarm/run-all.sh
```

# Launch remote tests

```
To be done
```