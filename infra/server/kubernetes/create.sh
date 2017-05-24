#!/bin/bash

# prerequisites for helm installed

# make sure you have enough quotas
# 
# for DB in cockroachdb mssql mysql oracle postgres;
# do
#     gcloud container clusters create $DB \
#         --zone=asia-southeast1-a --num-nodes=4 --preemptible \
#         --cluster-version=1.6.1 --async
# done;

# helm install --name cockroachdb stable/cockroachdb
# gcloud container clusters get-credentials cockroachdb --zone asia-southeast1-a --project playwithsql
# gcloud container clusters delete cockroachdb -q


gcloud container clusters create benchcluster \
    --zone=us-central1-f --num-nodes=4 --machine-type=n1-standard-1 \
    --cluster-version=1.6.2;

gcloud container clusters get-credentials benchcluster --zone us-central1-f --project playwithsql;