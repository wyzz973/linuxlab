#!/bin/bash
docker rmi envapp:v1 2>/dev/null
docker rm -f env-test 2>/dev/null
rm -rf /tmp/envapp
mkdir -p /tmp/envapp
