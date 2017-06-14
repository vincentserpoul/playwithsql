# playwithsql [![Go Report Card](https://goreportcard.com/badge/github.com/vincentserpoul/playwithsql)](https://goreportcard.com/report/github.com/vincentserpoul/playwithsql) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/a79310b8da354991a0b2b657a73f195f)](https://www.codacy.com/app/vincent_11/playwithsql?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=vincentserpoul/playwithsql&amp;utm_campaign=Badge_Grade)

implementing as immutable as possible data modelization and benchmarking it on different platforms.

## Disclaimer

the benchmark comparison is for very specific use case:
* Golang 1.8.3
* Containerized DBs, latest versions
* Specific schemas
* Used configurations
* GCP as cloud provider (or local, until docker for GCP is allowing experimental)
* n1-standard-1 as machine type

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
- [x] Reference hosted MySQL on GCP
- [x] Reference hosted Postgres on GCP
- [ ] Vitess
- [ ] Cloud Spanner

# Tested setups

- [x] local swarm cluster - local solo db container
- [ ] swarm cluster - gcloud solo db container
- [x] kubernetes cluster - gcloud solo db container
- [ ] swarm cluster - gcloud cluster db containers
- [ ] kubernetes cluster - gcloud cluster db containers

# Rebuilding the docker image

```
./build-docker-cmd-status.sh
```

# Launch tests

```
./status/run-test.sh
```

# Launch local status benches

```
./bench/status/swarm/run-all.sh (islatest or lateststatus or history)
```

# Launch remote tests on kubernetes

```
./bench/status/kubernetes/run-all.sh (islatest or lateststatus or history)
```

# TODO list

- [ ] Bench 1000000 loops and get the best of 5 runs for each db
- [ ] Test different storage (mounted standard disk, mounted ssd, local ssd)
- [ ] Test High Availability
- [ ] Leverage new go 1.8 capabilities (named queries, remove sqlx?)
- [ ] Vendor deps (github.com/kardianos/govendor)
