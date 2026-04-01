#!/bin/bash
docker pull nginx:latest 2>/dev/null
docker rm -f web-mount 2>/dev/null
rm -rf /tmp/web-content
