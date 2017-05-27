#!/bin/bash

pauseService () {
    gcloud sql instances patch gcppostgresbench --activation-policy NEVER -q;
}