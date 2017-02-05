#!/bin/bash

NODESCOUNT=$1
DOTOKEN=$2
DOSSHKEYPATH=$3
DOSSHFINGERPRINT=$4

#=========================
# Creating cluster members
#=========================
for n in $(seq 1 $NODESCOUNT); do
    docker-machine create --driver digitalocean --digitalocean-access-token=$DOTOKEN \
        --digitalocean-region=sgp1 \
        --digitalocean-image=ubuntu-16-10-x64 \
        --digitalocean-ssh-key-fingerprint=$DOSSHFINGERPRINT \
        --digitalocean-ssh-key-path=$DOSSHKEYPATH \
        node$n
done;

#===============
# Starting swarm
#===============
MANAGER_IP=$(docker-machine ip node1)

docker-machine ssh node1 docker swarm init --advertise-addr $MANAGER_IP

#===============
# Adding members
#===============
MANAGER_TOKEN=$(docker-machine ssh node1 docker swarm join-token --quiet manager)
WORKER_TOKEN=$(docker-machine ssh node1 docker swarm join-token --quiet worker)

for n in $(seq 2 $NODESCOUNT); do
    docker-machine ssh node$n docker swarm join --token $MANAGER_TOKEN $MANAGER_IP:2377
done;