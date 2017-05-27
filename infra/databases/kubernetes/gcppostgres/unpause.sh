#!/bin/bash

unpauseService () {
    gcloud sql instances patch gcppostgresbench --activation-policy ALWAYS -q;
}

unpauseService;