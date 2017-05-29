#!/bin/bash

initdb () {
    sleep 60s;
    POD_NAME=$(kubectl get po | grep mssql | awk '{ print $1 }');
    kubectl exec -i $POD_NAME -- /bin/bash -c 'echo "create database entityone_test;" > createdb.sql && /opt/mssql-tools/bin/sqlcmd -U sa -P thank5MsSQLforcingMe -i ./createdb.sql;';
    kubectl exec -i $POD_NAME -- /bin/bash -c 'echo "create database playwithsql;" > createdb.sql && /opt/mssql-tools/bin/sqlcmd -U sa -P thank5MsSQLforcingMe -i ./createdb.sql;';
}

removeService () {
    kubectl delete -f ./infra/databases/kubernetes/mssql/kube-solo.yml
}

runService () {
    removeService;
    kubectl create -f ./infra/databases/kubernetes/mssql/kube-solo.yml;
    while [ $(kubectl get po | grep Running | grep mssql | awk '{ print $3 }' | wc -l) -lt 1 ] ;do 
        sleep 1s;
    done;        
    initdb;
}

runService;