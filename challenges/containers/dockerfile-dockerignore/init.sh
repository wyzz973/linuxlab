#!/bin/bash
docker rmi ignore-demo:v1 2>/dev/null
rm -rf /tmp/ignore-demo
mkdir -p /tmp/ignore-demo/logs
