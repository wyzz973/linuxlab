#!/bin/bash
docker volume create shared-data
docker run -d --name producer -v shared-data:/shared alpine sh -c 'while true; do date > /shared/timestamp.txt; sleep 1; done'
docker run -d --name consumer -v shared-data:/shared alpine sleep 3600
sleep 3
docker exec consumer cat /shared/timestamp.txt > /tmp/shared-data.txt
