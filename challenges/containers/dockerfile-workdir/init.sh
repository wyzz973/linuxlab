#!/bin/bash
docker rmi workdir-app:v1 2>/dev/null
rm -rf /tmp/workdir-app
mkdir -p /tmp/workdir-app
