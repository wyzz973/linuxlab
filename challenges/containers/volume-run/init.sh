#!/bin/bash
docker pull alpine:latest 2>/dev/null
docker rm -f writer reader 2>/dev/null
docker volume rm persist-data 2>/dev/null
rm -f /tmp/persist-check.txt
