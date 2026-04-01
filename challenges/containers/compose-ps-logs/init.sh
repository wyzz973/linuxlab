#!/bin/bash
docker pull nginx:alpine 2>/dev/null
docker rm -f logs-web 2>/dev/null
rm -f /tmp/compose-web-logs.txt /tmp/compose-ps.txt
