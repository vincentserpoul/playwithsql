#!/bin/bash

# install docker compose

curl -L https://github.com/docker/compose/releases/download/1.9.0/docker-compose-`uname -s`-`uname -m` > docker-compose
chmod +x /opt/bin/docker-compose
sudo mkdir -p /opt/bin
sudo mv docker-compose /opt/bin
sudo chmod +x /opt/bin/docker-compose

rm -rf playwithsql
git clone https://github.com/vincentserpoul/playwithsql.git