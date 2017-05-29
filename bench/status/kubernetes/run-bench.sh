#!/bin/bash

export DB=$1;
export HOST=$2;
export SCH=$3;
export LOOPS=$4;
export MAXCONNS=$5;

envsubst < ./bench/status/kubernetes/kube-bench.yml | kubectl delete -f -;
envsubst < ./bench/status/kubernetes/kube-bench.yml | kubectl create -f -;

while [ $(kubectl get po -a | grep bench | grep Completed | grep $DB | awk '{ print $3 }' | wc -l) -lt 1 ] ;do 
    sleep 1s;
done;

POD_NAME=$(kubectl get po -a | grep bench | grep Completed | grep $DB | awk '{ print $1 }');
kubectl logs $POD_NAME >> ./bench/status/kubernetes/$SCH/results.log;
echo ","  >> ./bench/status/kubernetes/$SCH/results.log;
kubectl delete -f ./infra/databases/kubernetes/$DB/kube-solo.yml;