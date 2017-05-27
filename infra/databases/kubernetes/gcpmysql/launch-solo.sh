#!/bin/bash

removeService () {
    kubectl delete -f ./infra/databases/kubernetes/gcpmysql/kube-solo.yml;
}

runService () {
    removeService;
    kubectl create secret generic cloudsql-instance-credentials --from-file=credentials.json=./infra/databases/kubernetes/gcpmysql/credentials.json;
    kubectl create secret generic cloudsql-db-credentials --from-literal=username=root --from-literal=password=test;    
    kubectl create -f ./infra/databases/kubernetes/gcpmysql/kube-solo.yml;
}

runService;