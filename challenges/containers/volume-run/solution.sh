#!/bin/bash
docker volume create persist-data
docker run --name writer -v persist-data:/data alpine sh -c 'echo "data survives" > /data/message.txt'
docker rm writer
docker run --rm --name reader -v persist-data:/data alpine cat /data/message.txt > /tmp/persist-check.txt
