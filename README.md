[![Go Report Card](https://goreportcard.com/badge/github.com/vincentserpoul/playwithsql)](https://goreportcard.com/report/github.com/vincentserpoul/playwithsql)
[![codebeat badge](https://codebeat.co/badges/df4fb8c7-3472-46ff-a9c8-5fa72008269c)](https://codebeat.co/projects/github-com-vincentserpoul-playwithsql-master)

# playwithsql

implementing as immutable as possible data modelization and benchmarking it on different platforms.

## Disclaimer

the benchmark comparison is for very specific use case:
* Golang 1.8.1
* Containerized DBs
* Specific schemas
* Used configurations
* GCP as cloud provider

Hence, they can't be used to affirm that this or this db is better.
*The context matters!*

## Results

Just check them [here](https://playwithsql-summary.surge.sh)!

# Tested dbs

- [x] MySQL
- [x] MariaDB
- [x] PerconaDB
- [x] Postgres
- [x] CockroachDB
- [x] SQLite
- [x] Microsoft SQL Server
- [x] Oracle
- [ ] Vitess
- [ ] Reference hosted MySQL on GCP
- [ ] Cloud Spanner

# Tested setups

- [x] local swarm cluster - local solo db container
- [ ] swarm cluster - gcloud solo db container
- [x] kubernetes cluster - gcloud solo db container
- [ ] swarm cluster - gcloud cluster db containers
- [ ] kubernetes cluster - gcloud cluster db containers

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

# Launch remote tests on kubernetes

```
./bench/status/kubernetes/run-all.sh
```

# TODO list

- [ ] Leverage new go 1.8.1 capabilities (context, timeouts done, left with named queries, remove sqlx?)
- [ ] Vendor deps (github.com/kardianos/govendor)
- [ ] Test different storage
- [ ] Test High Availability
