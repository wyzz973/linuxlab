#!/bin/bash
docker network create app-network
docker run -d --network app-network --name web-app nginx
docker run -d --network app-network --name test-client alpine sleep 3600
sleep 2
docker exec test-client ping -c 3 web-app > /tmp/ping-result.txt 2>&1
