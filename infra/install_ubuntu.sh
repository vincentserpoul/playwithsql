#!/bin/bash

## sqlite3 needs
apt-get install -y sqlite3 gcc



## ORACLE NEEDS
apt-get install -y alien libaio1 unzip pkg-config
mkdir -p /opt/oracle
cd /opt/oracle 
unzip /root/instantclient-basic-linux.x64-12.1.0.2.0.zip
unzip /root/instantclient-sdk-linux.x64-12.1.0.2.0.zip
cd /opt/oracle/instantclient_12_1
echo "export LD_LIBRARY_PATH=/opt/oracle/instantclient_12_1:/opt/oracle/instantclient_12_1/sdk/include" >> ~/.bashrc 
echo "export PKG_CONFIG_PATH=/opt/oracle" >> ~/.bashrc 
echo "export PATH=/opt/oracle/instantclient_12_1:$PATH" >> ~/.bashrc 
echo "export ORACLE_HOME=/opt/oracle/instantclient_12_1:/opt/oracle/instantclient_12_1/sdk/include" >> ~/.bashrc 
exec bash 
ln -s /opt/oracle/libclntsh.so.12.1 /opt/oracle/libclntsh.so
ln -s /opt/oracle/libocci.so.12.1 /opt/oracle/libocci.so

## Install GO
cd /root
wget https://storage.googleapis.com/golang/go1.7.4.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.7.4.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc 
echo "export GOPATH=$HOME" >> ~/.bashrc 
exec bash

## Install playwithsql
mkdir -p $GOPATH/src/github.com/vincentserpoul
rm -rf $GOPATH/src/github.com/vincentserpoul/playwithsql
git clone https://github.com/vincentserpoul/playwithsql.git $GOPATH/src/github.com/vincentserpoul/playwithsql
cd $GOPATH/src/github.com/vincentserpoul/playwithsql/
cp -aL $GOPATH/src/github.com/vincentserpoul/playwithsql/infra/oci8.pc $PKG_CONFIG_PATH/oci8.pc
go install ./...