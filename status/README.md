# bench golang

on Intel(R) Core(TM) i7-5500U CPU @ 2.40GHz, 16GB ram

```
$ ./status/run-test.sh                                                                              

====================

Cockroachdb

islatest
--------
BenchmarkCreate-4                    	     500	  10468330 ns/op
BenchmarkUpdateStatus-4              	     500	   9011778 ns/op
BenchmarkSelectEntityoneByStatus-4   	     300	  18448235 ns/op
BenchmarkSelectEntityoneOneByPK-4    	     300	  17290208 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	26.662s

lateststatus
------------
BenchmarkCreate-4                    	     300	  13330483 ns/op
BenchmarkUpdateStatus-4              	     500	  10548141 ns/op
BenchmarkSelectEntityoneByStatus-4   	     300	  15148713 ns/op
BenchmarkSelectEntityoneOneByPK-4    	     500	   9512001 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	24.508s

history
-------
BenchmarkCreate-4                    	     500	  10360480 ns/op
BenchmarkUpdateStatus-4              	     500	  10700921 ns/op
BenchmarkSelectEntityoneByStatus-4   	    3000	   1349362 ns/op
BenchmarkSelectEntityoneOneByPK-4    	    5000	   1199731 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	23.754s

====================

MSSQL

islatest
--------
BenchmarkCreate-4                    	    2000	   2151683 ns/op
BenchmarkUpdateStatus-4              	    3000	   1502009 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    410183 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    371751 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	17.162s

lateststatus
------------
BenchmarkCreate-4                    	    2000	   2441415 ns/op
BenchmarkUpdateStatus-4              	    2000	   1934141 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    409696 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    358960 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	17.056s

history
-------
BenchmarkCreate-4                    	    2000	   1877509 ns/op
BenchmarkUpdateStatus-4              	    3000	   1361104 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    358764 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    324062 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	15.149s

====================

MySQL 

islatest
--------
BenchmarkCreate-4                    	    2000	   2499545 ns/op
BenchmarkUpdateStatus-4              	    2000	   2526334 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    426611 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    418569 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	19.206s

lateststatus
------------
BenchmarkCreate-4                    	    2000	   2069397 ns/op
BenchmarkUpdateStatus-4              	    2000	   2056994 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    422564 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    368221 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	16.790s

history
-------
BenchmarkCreate-4                    	    2000	   1740044 ns/op
BenchmarkUpdateStatus-4              	    2000	   2054322 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    368390 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    319086 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	15.055s

====================

Postgres

islatest
--------
BenchmarkCreate-4                    	    2000	   2063406 ns/op
BenchmarkUpdateStatus-4              	    2000	   1775458 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    634958 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    491673 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	19.600s

lateststatus
------------
BenchmarkCreate-4                    	    3000	   1724629 ns/op
BenchmarkUpdateStatus-4              	    2000	   1644451 ns/op
BenchmarkSelectEntityoneByStatus-4   	    5000	    754812 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    494935 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	17.760s

history
-------
BenchmarkCreate-4                    	    3000	   1751059 ns/op
BenchmarkUpdateStatus-4              	    2000	   2095055 ns/op
BenchmarkSelectEntityoneByStatus-4   	    5000	    639743 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    313535 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	16.323s

====================

SQLite

islatest
--------
BenchmarkCreate-4                    	    1000	   5139560 ns/op
BenchmarkUpdateStatus-4              	    1000	   5128113 ns/op
BenchmarkSelectEntityoneByStatus-4   	   50000	     93758 ns/op
BenchmarkSelectEntityoneOneByPK-4    	  100000	     90393 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	26.589s

lateststatus
------------
BenchmarkCreate-4                    	    1000	   6139074 ns/op
BenchmarkUpdateStatus-4              	    1000	   5803292 ns/op
BenchmarkSelectEntityoneByStatus-4   	   50000	    124262 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   50000	    114864 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	27.136s

history
-------
BenchmarkCreate-4                    	    1000	   5519332 ns/op
BenchmarkUpdateStatus-4              	    1000	   5600910 ns/op
BenchmarkSelectEntityoneByStatus-4   	   50000	     91100 ns/op
BenchmarkSelectEntityoneOneByPK-4    	  100000	     77634 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	26.121s


```