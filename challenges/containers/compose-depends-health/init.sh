#!/bin/bash
docker pull redis:7-alpine nginx:alpine 2>/dev/null
docker rm -f health-db health-web 2>/dev/null
rm -rf /tmp/compose-health
mkdir -p /tmp/compose-health
