#!/bin/bash
docker pull nginx:alpine 2>/dev/null
docker rm -f vol-web 2>/dev/null
docker volume rm web-data 2>/dev/null
rm -rf /tmp/compose-vol
mkdir -p /tmp/compose-vol/config
