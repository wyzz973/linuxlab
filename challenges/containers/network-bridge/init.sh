#!/bin/bash
docker pull alpine:latest 2>/dev/null
docker rm -f bridge-test-1 bridge-test-2 2>/dev/null
rm -f /tmp/bridge-info.txt /tmp/container-ips.txt
