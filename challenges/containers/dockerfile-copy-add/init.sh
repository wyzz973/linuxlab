#!/bin/bash
docker rmi webapp:v1 2>/dev/null
rm -rf /tmp/webapp
mkdir -p /tmp/webapp
