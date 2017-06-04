# DBs max connections

# CockroachDB

according to their [example repo](https://github.com/cockroachdb/examples-go), it should be 2x number of CPUs + 1.

# MSSQL

[It seems to be 32,767](https://docs.microsoft.com/en-us/sql/sql-server/maximum-capacity-specifications-for-sql-server#Engine)

# MariaDB, MySQL, Percona

Determined by the arg --max-connections

# Postgres

```
# cat /var/lib/postgresql/data/postgresql.conf | grep max_connections
max_connections = 100			# (change requires restart)
```

```
SHOW max_connections;
```

