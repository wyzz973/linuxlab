#!/bin/bash
docker pull nginx:latest alpine:latest 2>/dev/null
docker rm -f proxy-server app-server 2>/dev/null
docker network rm frontend backend 2>/dev/null
