#!/bin/bash

initdb () {
    sleep 80s;
    POD_NAME=$(kubectl get po | grep percona | awk '{ print $1 }');
    kubectl exec -i $POD_NAME -- mysql -u root -ptest -e 'CREATE DATABASE entityone_test';
    kubectl exec -i $POD_NAME -- mysql -u root -ptest -e 'CREATE DATABASE playwithsql';
}

removeService () {
    kubectl delete -f ./infra/databases/kubernetes/percona/kube-solo.yml
}

runService () {
    removeService;
    kubectl create -f ./infra/databases/kubernetes/percona/kube-solo.yml;
    while [ $(kubectl get po | grep Running | grep percona | awk '{ print $3 }' | wc -l) -lt 1 ] ;do 
        sleep 1s;
    done;      
    initdb;
}

runService;