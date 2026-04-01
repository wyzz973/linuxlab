#!/bin/bash
docker rmi labeled-app:v2 2>/dev/null
rm -rf /tmp/labeled-app
mkdir -p /tmp/labeled-app
