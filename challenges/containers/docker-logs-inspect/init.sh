#!/bin/bash
docker pull nginx:latest 2>/dev/null
docker rm -f log-demo 2>/dev/null
rm -f /tmp/container-logs.txt /tmp/container-ip.txt
