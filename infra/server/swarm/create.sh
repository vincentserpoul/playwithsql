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
gcloud deployment-manager deployments create docker \
    --config https://download.docker.com/gcp/edge/Docker.jinja \
    --properties managerCount:3,workerCount:1,managerMachineType:n1-standard-1,workerMachineType:n1-standard-1,managerDiskType:pd-ssd,workerDiskType:pd-ssd

# ssh tunneling enable
# gcloud compute ssh --project playwithsql --zone us-central1-f docker-manager-1 -- -NL localhost:2374:/var/run/docker.sock &
# export DOCKER_HOST=localhost:2374