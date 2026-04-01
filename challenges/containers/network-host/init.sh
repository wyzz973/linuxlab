#!/bin/bash
docker pull nginx:latest 2>/dev/null
docker rm -f host-nginx 2>/dev/null
rm -f /tmp/host-network-info.txt
