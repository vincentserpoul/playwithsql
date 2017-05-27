#!/bin/bash

pauseService () {
    gcloud sql instances patch gcpmysqlbench --activation-policy NEVER -q;
}

pauseService;