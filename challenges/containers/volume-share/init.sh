#!/bin/bash
docker pull alpine:latest 2>/dev/null
docker rm -f producer consumer 2>/dev/null
docker volume rm shared-data 2>/dev/null
rm -f /tmp/shared-data.txt
