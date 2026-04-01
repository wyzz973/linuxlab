#!/bin/bash
docker pull nginx:alpine 2>/dev/null
cd /tmp && docker compose -f /Users/sd3/Desktop/project/challenges/containers/compose-basics/docker-compose.yaml down 2>/dev/null
docker rm -f compose-web 2>/dev/null
