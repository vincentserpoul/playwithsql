#!/bin/bash

# install docker compose

curl -L https://github.com/docker/compose/releases/download/1.9.0/docker-compose-`uname -s`-`uname -m` > docker-compose
chmod +x docker-compose
sudo mkdir -p /opt/bin
# sudo cp docker-compose docker-compose.sav
sudo mv docker-compose /opt/bin

docker-compose version

rm -rf playwithsql
git clone https://github.com/vincentserpoul/playwithsql.git