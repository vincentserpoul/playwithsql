# DBs max connections

# CockroachDB

according to their [example repo](https://github.com/cockroachdb/examples-go), it should be 2x number of CPUs + 1.

# MSSQL

[It seems to be 32,767](https://docs.microsoft.com/en-us/sql/sql-server/maximum-capacity-specifications-for-sql-server#Engine)

# MariaDB, MySQL, Percona

Determined by the arg --max-connections

# Oracle

in SQLPlus

```
SQL> SET LIN 200
SQL> select * from v$resource_limit where resource_name in ('processes', 'sessions', 'transactions');

RESOURCE_NAME		       CURRENT_UTILIZATION MAX_UTILIZATION INITIAL_ALLOCATION			    LIMIT_VALUE
------------------------------ ------------------- --------------- ---------------------------------------- ----------------------------------------
processes					29		90	  100					   100
sessions					30		72	  172					   172
transactions					 0		 0	  189				     UNLIMITED
```

# Postgres

```
# cat /var/lib/postgresql/data/postgresql.conf | grep max_connections
max_connections = 100			# (change requires restart)
```

```
SHOW max_connections;
```

