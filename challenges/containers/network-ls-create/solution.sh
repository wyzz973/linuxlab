#!/bin/bash
docker network ls > /tmp/network-list.txt
docker network create my-network
docker network create --subnet=192.168.100.0/24 isolated-net
docker network ls > /tmp/network-list-after.txt
