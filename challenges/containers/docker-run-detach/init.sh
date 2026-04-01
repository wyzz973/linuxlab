#!/bin/bash
docker pull nginx:latest 2>/dev/null
docker rm -f my-nginx 2>/dev/null
