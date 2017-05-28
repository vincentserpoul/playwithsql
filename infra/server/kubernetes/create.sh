#!/bin/bash

gcloud container clusters create benchcluster \
    --zone=us-central1-f --num-nodes=4 --machine-type=n1-standard-1 \
    --cluster-version=1.6.4;

gcloud container clusters get-credentials benchcluster --zone us-central1-f --project playwithsql;

# gcloud container clusters delete benchcluster --zone us-central1-f -q