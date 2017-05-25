#!/bin/bash

# PREPARATION
# You need to create the credentials.json, follow this tutorial
# https://cloud.google.com/sql/docs/mysql/connect-container-engine
# place your credentials in ./infra/databases/kubernetes/gcpmysqlbench/credentials.json

initdb () {
    # still not functional, to do manually on https://console.cloud.google.com/sql/instances/gcpmysqlbench/databases?project=playwithsql
    ACCESS_TOKEN=$(GOOGLE_APPLICATION_CREDENTIALS="./infra/databases/kubernetes/gcpmysql/credentials.json" gcloud auth application-default print-access-token);
    curl --header "Authorization: Bearer ${ACCESS_TOKEN}" \
        https://www.googleapis.com/sql/v1beta4/projects/playwithsql/instances/playwithsql:us-central1:gcpmysqlbench/databases/playwithsql -X DELETE;    
    curl --header "Authorization: Bearer ${ACCESS_TOKEN}" \
        --header 'Content-Type: application/json' \
        --data '{"project": "playwithsql", "instance": "playwithsql:us-central1:gcpmysqlbench", "name": "playwithsql"}' \
        https://www.googleapis.com/sql/v1beta4/projects/playwithsql/instances/playwithsql:us-central1:gcpmysqlbench/databases -X POST
}

removeService () {
    gcloud beta sql instances delete gcpmysqlbench -q;
}

runService () {
    # as you cannot reuse an instance name for up to a week after you have deleted an instance, we better not do that
    # removeService;
    # gcloud beta sql instances create gcpmysqlbench --tier=db-n1-standard-1 --region=us-central1 --database-version MYSQL_5_7;
    # gcloud beta sql users set-password root % --instance gcpmysqlbench --password test;
    # kubectl create secret generic cloudsql-instance-credentials --from-file=credentials.json=./infra/databases/kubernetes/gcpmysql/credentials.json;
    # kubectl create secret generic cloudsql-db-credentials --from-literal=username=root --from-literal=password=test;
    # kubectl create -f ./infra/databases/kubernetes/gcpmysql/cloud-proxy.yml;
    read -n 1 -p "Create the database directly on the UI (https://console.cloud.google.com/sql/instances/gcpmysqlbench/databases?project=playwithsql) and press enter to continue.";
    initdb;
}

runService;