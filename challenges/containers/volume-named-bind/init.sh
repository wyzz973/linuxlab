#!/bin/bash
docker pull alpine:latest 2>/dev/null
docker rm -f vol-container-1 vol-container-2 2>/dev/null
docker volume rm named-vol 2>/dev/null
rm -rf /tmp/bind-data
