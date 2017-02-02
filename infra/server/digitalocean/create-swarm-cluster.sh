#!/bin/bash

DOTOKEN=$1
DOSSHKEYPATH=$2

docker-machine create --driver digitalocean --digitalocean-access-token=$DOTOKEN --digitalocean-region=sgp1 --digitalocean-image=coreos-stable --digitalocean-ssh-key-fingerprint=9a:03:08:b9:bc:a0:fc:6a:20:6d:86:9f:18:c0:e8:ed --digitalocean-ssh-user=core --digitalocean-ssh-key-path=$DOSSHKEYPATH --swarm node1