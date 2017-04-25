#!/bin/bash

# Before running

# set your default project
# gcloud config set project playwithsql

# activate Google Cloud RuntimeConfig API
# https://console.developers.google.com/apis/api/runtimeconfig.googleapis.com/overview?project=playwithsql

# Enable container registry
# https://console.cloud.google.com/gcr/images/playwithsql?project=playwithsql

# https://beta.docker.com/docs/gcp/

# Create docker deployment
gcloud deployment-manager deployments create docker --config https://download.docker.com/gcp/edge/Docker.jinja --properties managerCount:3,workerCount:1,zone:asia-east1-c