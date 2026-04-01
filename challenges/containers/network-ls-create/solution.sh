#!/bin/bash
docker network ls > /tmp/network-list.txt
docker network create my-network
docker network create --subnet=172.28.0.0/16 isolated-net
docker network ls > /tmp/network-list-after.txt
