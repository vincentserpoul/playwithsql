#!/bin/bash

DOTOKEN=$1
DOSSHKEYPATH=$2
DOSSHFINGERPRINT=$3

#=========================
# Creating cluster members
#=========================
for n in $(seq 1 3); do
    docker-machine create --driver digitalocean --digitalocean-access-token=$DOTOKEN \
        --digitalocean-region=sgp1 \
        --digitalocean-image=coreos-stable \
        --digitalocean-ssh-key-fingerprint=$DOSSHFINGERPRINT \
        --digitalocean-ssh-user=core \
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

for n in $(seq 2 3); do
    docker-machine ssh node$n docker swarm join --token $MANAGER_TOKEN $MANAGER_IP:2377
done;