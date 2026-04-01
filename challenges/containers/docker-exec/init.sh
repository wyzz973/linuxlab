#!/bin/bash
docker pull nginx:latest 2>/dev/null
docker rm -f exec-demo 2>/dev/null
rm -f /tmp/container-processes.txt
