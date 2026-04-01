#!/bin/bash
mkdir -p /tmp/bind-data
docker volume create named-vol
docker run --name vol-container-1 -v named-vol:/data alpine sh -c 'echo "named volume data" > /data/test.txt'
docker run --name vol-container-2 -v /tmp/bind-data:/data alpine sh -c 'echo "bind mount data" > /data/test.txt'
