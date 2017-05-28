# bench golang

on Intel(R) Core(TM) i7-5500U CPU @ 2.40GHz, 16GB ram

```
$ ./status/run-test.sh                                                                              

cockroachdb

islatest
BenchmarkCreate-4                    	     500	  12137799 ns/op
BenchmarkUpdateStatus-4              	     500	  10547835 ns/op
BenchmarkSelectEntityoneByStatus-4   	     200	  28260754 ns/op
BenchmarkSelectEntityoneOneByPK-4    	     200	  25584838 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	30.626s

lateststatus
BenchmarkCreate-4                    	     500	  11064733 ns/op
BenchmarkUpdateStatus-4              	     500	  10752911 ns/op
BenchmarkSelectEntityoneByStatus-4   	     200	  24842710 ns/op
BenchmarkSelectEntityoneOneByPK-4    	    1000	   7142678 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	29.171s


mssql

islatest
BenchmarkCreate-4                    	    2000	   2194208 ns/op
BenchmarkUpdateStatus-4              	    3000	   1493990 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    408655 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    370540 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	17.208s

lateststatus
BenchmarkCreate-4                    	    2000	   2109077 ns/op
BenchmarkUpdateStatus-4              	    3000	   1734890 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    385789 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    339503 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	17.227s


mysql

islatest
BenchmarkCreate-4                    	    2000	   2047932 ns/op
BenchmarkUpdateStatus-4              	    2000	   2099221 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    377041 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    323345 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	15.899s

lateststatus
BenchmarkCreate-4                    	    2000	   1811794 ns/op
BenchmarkUpdateStatus-4              	    2000	   1965645 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    370748 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    320367 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	15.019s


oracle

islatest
BenchmarkCreate-4                    	    2000	   2212492 ns/op
BenchmarkUpdateStatus-4              	    2000	   2077987 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    489922 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    968997 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	24.098s

lateststatus
BenchmarkCreate-4                    	    2000	   2243791 ns/op
BenchmarkUpdateStatus-4              	    2000	   2188894 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    455243 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    420578 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	18.936s


postgres

islatest
BenchmarkCreate-4                    	    3000	   1639382 ns/op
BenchmarkUpdateStatus-4              	    3000	   1515541 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    476851 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    379906 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	18.612s

lateststatus
BenchmarkCreate-4                    	    3000	   1667420 ns/op
BenchmarkUpdateStatus-4              	    3000	   1607581 ns/op
BenchmarkSelectEntityoneByStatus-4   	   10000	    643105 ns/op
BenchmarkSelectEntityoneOneByPK-4    	   10000	    372732 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	20.491s


SQLite

islatest
BenchmarkCreate-4                    	    1000	   6062866 ns/op
BenchmarkUpdateStatus-4              	    1000	   6156758 ns/op
BenchmarkSelectEntityoneByStatus-4   	   50000	    103742 ns/op
BenchmarkSelectEntityoneOneByPK-4    	  100000	     94753 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	29.896s

lateststatus
BenchmarkCreate-4                    	    1000	   6479264 ns/op
BenchmarkUpdateStatus-4              	    1000	   6338754 ns/op
BenchmarkSelectEntityoneByStatus-4   	   50000	    101534 ns/op
BenchmarkSelectEntityoneOneByPK-4    	  100000	     85954 ns/op
PASS
ok  	github.com/vincentserpoul/playwithsql/status	29.515s
```