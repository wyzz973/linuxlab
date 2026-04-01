#!/bin/bash
docker pull nginx:latest 2>/dev/null
docker rm -f cp-demo 2>/dev/null
rm -f /tmp/upload.txt /tmp/nginx.conf
