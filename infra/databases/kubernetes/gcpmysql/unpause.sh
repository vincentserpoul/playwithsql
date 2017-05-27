#!/bin/bash

unpauseService () {
    gcloud sql instances patch gcpmysqlbench --activation-policy ALWAYS -q;
}

unpauseService;