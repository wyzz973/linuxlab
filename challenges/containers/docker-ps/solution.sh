#!/bin/bash
docker run -d --name running-web nginx
docker run --name stopped-task alpine echo "done"
docker ps -a > /tmp/all-containers.txt
docker ps > /tmp/running-containers.txt
