#!/bin/bash


cd ./KVStore/Server
go build Server.go && mv Server ../../bin/
cd -
cd ./KVStore/Client
go build Client.go && mv Client ../../bin/
cd -
