#!/bin/bash
docker pull nginx:latest 2>/dev/null
docker rm -f web-server 2>/dev/null
