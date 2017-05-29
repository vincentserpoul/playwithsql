#!/bin/bash

initdb () {
    sleep 60s;
    POD_NAME=$(kubectl get po | grep cockroach | awk '{ print $1 }');
    kubectl exec -i $POD_NAME -- ./cockroach sql --insecure \
        --execute="CREATE DATABASE entityone_test;";
    kubectl exec -i $POD_NAME -- ./cockroach sql --insecure \
        --execute="CREATE DATABASE playwithsql;";
}

removeService () {
    kubectl delete -f ./infra/databases/kubernetes/cockroachdb/kube-solo.yml
}

runService () {
    removeService;
    kubectl create -f ./infra/databases/kubernetes/cockroachdb/kube-solo.yml;
    while [ $(kubectl get po | grep Running | grep cockroachdb | awk '{ print $3 }' | wc -l) -lt 1 ] ;do 
        sleep 1s;
    done;       
    initdb;
}

runService;