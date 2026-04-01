#!/bin/bash
docker pull alpine:latest 2>/dev/null
docker rm -f temp-container 2>/dev/null
rm -f /tmp/container-test.txt
