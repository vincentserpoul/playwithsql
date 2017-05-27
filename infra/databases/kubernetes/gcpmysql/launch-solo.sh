#!/bin/bash

removeService () {
    kubectl delete -f ./infra/databases/kubernetes/gcpmysql/kube-solo.yml;
}

runService () {
    removeService;
    kubectl create -f ./infra/databases/kubernetes/gcpmysql/kube-solo.yml;
}

runService;