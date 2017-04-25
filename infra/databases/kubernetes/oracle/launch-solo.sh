#!/bin/bash

initdb () {
    sleep 45s;
    POD_NAME=$(kubectl get po | grep oracle | awk '{ print $1 }');
    kubectl exec -i $POD_NAME -- /bin/bash -c 'ORACLE_HOME="/u01/app/oracle/product/11.2.0/xe" ORACLE_SID="XE" u01/app/oracle/product/11.2.0/xe/bin/sqlplus -s /nolog <<EOF
connect system/oracle
    create user playwithsql identified by "dev";
    grant all privileges to playwithsql;
    create user entityone_test identified by "dev";
    grant all privileges to entityone_test;
quit
EOF';
}

removeService () {
    kubectl delete -f ./infra/databases/kubernetes/oracle/kube-solo.yml
}

runService () {
    removeService;
    kubectl create -f ./infra/databases/kubernetes/oracle/kube-solo.yml;
    while [ $(kubectl get po | grep Running | grep oracle | awk '{ print $3 }' | wc -l) -lt 1 ] ;do 
        sleep 1s;
    done;     
    initdb;
}

runService;