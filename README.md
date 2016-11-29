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

1x coreos 16 GB Memory / 160 GB Disk / SGP1 - CoreOS 1185.3.0 (stable)
1x ubuntu 16 GB Memory / 160 GB Disk / SGP1 - Ubuntu 16.04.1 x64

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

# The results

## SQLite3

```
go test -db=sqlite -bench=.  -test.benchtime=60s;rm -f ./test.db;
```

BenchmarkCreate-8                       	   10000	   7518683 ns/op
BenchmarkUpdateStatus-8                 	   10000	   7913274 ns/op
BenchmarkSelectEntityoneOneByStatus-8   	  500000	    205393 ns/op
BenchmarkSelectEntityoneOneByPK-8       	  500000	    182951 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status/islatest	354.145s

## CockroachDB - 1 container

BenchmarkCreate-8                       	    5000	  15425321 ns/op
BenchmarkUpdateStatus-8                 	   10000	  11939611 ns/op
BenchmarkSelectEntityoneOneByStatus-8   	     300	 222693116 ns/op
BenchmarkSelectEntityoneOneByPK-8       	       3	28799012640 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status/islatest	469.577s

## CockroachDB - 5 containers

BenchmarkCreate-8                       	    5000	  21197283 ns/op
BenchmarkUpdateStatus-8                 	    5000	  25218792 ns/op
BenchmarkSelectEntityoneOneByStatus-8   	     500	 224125369 ns/op
BenchmarkSelectEntityoneOneByPK-8       	       2	31793801784 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status/islatest	467.998s

## PerconaDB 5.7.15

BenchmarkCreate-8                       	   10000	   6687659 ns/op
BenchmarkUpdateStatus-8                 	   20000	   4738893 ns/op
BenchmarkSelectEntityoneOneByStatus-8   	   50000	   1844723 ns/op
BenchmarkSelectEntityoneOneByPK-8       	   50000	   1706539 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status/islatest	422.740s

## Postgres 9.6.1

BenchmarkCreate-8                       	   20000	   5760065 ns/op
BenchmarkUpdateStatus-8                 	   20000	   4469207 ns/op
BenchmarkSelectEntityoneOneByStatus-8   	   50000	   2108753 ns/op
BenchmarkSelectEntityoneOneByPK-8       	   50000	   1974964 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status/islatest	562.721s