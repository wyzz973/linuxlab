#!/bin/bash
cd "$(dirname "$0")"
docker compose up -d
sleep 2
curl -s http://localhost:8092 > /dev/null 2>&1
docker compose logs web > /tmp/compose-web-logs.txt 2>&1
docker compose ps > /tmp/compose-ps.txt
