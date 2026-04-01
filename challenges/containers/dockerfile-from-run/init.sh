#!/bin/bash
docker rmi myapp:v1 2>/dev/null
rm -rf /tmp/myapp
mkdir -p /tmp/myapp
