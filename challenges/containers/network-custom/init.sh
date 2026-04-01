#!/bin/bash
docker pull nginx:latest alpine:latest 2>/dev/null
docker rm -f web-app test-client 2>/dev/null
docker network rm app-network 2>/dev/null
rm -f /tmp/ping-result.txt
