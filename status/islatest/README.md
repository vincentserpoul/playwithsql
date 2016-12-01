# The db schema

![schema](status_islatest.png)

# The setup

* 2 cores / 4 GB Memory / 60 GB Disk / SGP1 - CoreOS 1185.3.0 (stable)
* 2 cores / 4 GB Memory / 60 GB Disk / SGP1 - Ubuntu 16.04.1 x64

# The results (ordered by overall time asc)

## SQLite3

```
BenchmarkCreate-2                       	   10000	   6983894 ns/op
BenchmarkUpdateStatus-2                 	   10000	   7501474 ns/op
BenchmarkSelectEntityoneOneByStatus-2   	  500000	    208845 ns/op
BenchmarkSelectEntityoneOneByPK-2       	  500000	    167025 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status/islatest	338.031s
```

## MariaDB 10.1.19

```
BenchmarkCreate-2                       	   10000	   6013189 ns/op
BenchmarkUpdateStatus-2                 	   20000	   4569070 ns/op
BenchmarkSelectEntityoneOneByStatus-2   	   50000	   1339161 ns/op
BenchmarkSelectEntityoneOneByPK-2       	  100000	   1277267 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status/islatest	421.562s
```

## PerconaDB 5.7.15

```
BenchmarkCreate-2                       	   10000	   6686537 ns/op
BenchmarkUpdateStatus-2                 	   20000	   4794422 ns/op
BenchmarkSelectEntityoneOneByStatus-2   	   50000	   1883283 ns/op
BenchmarkSelectEntityoneOneByPK-2       	   50000	   1721210 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status/islatest	426.369s
```

## MySQL 8.0.0

```
BenchmarkCreate-2                       	   10000	   6516587 ns/op
BenchmarkUpdateStatus-2                 	   20000	   5053199 ns/op
BenchmarkSelectEntityoneOneByStatus-2   	   50000	   1872660 ns/op
BenchmarkSelectEntityoneOneByPK-2       	   50000	   1663585 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status/islatest	427.266s
```

## Postgres 9.6.1

```
BenchmarkCreate-2                       	   10000	   6127597 ns/op
BenchmarkUpdateStatus-2                 	   20000	   4521845 ns/op
BenchmarkSelectEntityoneOneByStatus-2   	   50000	   2212562 ns/op
BenchmarkSelectEntityoneOneByPK-2       	   50000	   1871948 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status/islatest	443.242s
```

## CockroachDB - 5 containers

```
BenchmarkCreate-2                       	    2000	  68236742 ns/op
BenchmarkUpdateStatus-2                 	    2000	  52848647 ns/op
BenchmarkSelectEntityoneOneByStatus-2   	     500	 173088706 ns/op
BenchmarkSelectEntityoneOneByPK-2       	      10	7981009086 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status/islatest	454.705s
```

## CockroachDB - 1 container

```
BenchmarkCreate-2                       	    5000	  19749638 ns/op
BenchmarkUpdateStatus-2                 	    5000	  14632276 ns/op
BenchmarkSelectEntityoneOneByStatus-2   	     500	 196795986 ns/op
BenchmarkSelectEntityoneOneByPK-2       	       3	29519557994 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status/islatest	477.142s
```