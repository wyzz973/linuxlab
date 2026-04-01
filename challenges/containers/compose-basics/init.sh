#!/bin/bash
docker pull nginx:alpine 2>/dev/null
cd "$(dirname "$0")" && docker compose down 2>/dev/null
docker rm -f compose-web 2>/dev/null
