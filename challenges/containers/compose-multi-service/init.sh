#!/bin/bash
docker pull nginx:alpine node:18-alpine redis:7-alpine 2>/dev/null
rm -rf /tmp/multi-app
mkdir -p /tmp/multi-app
