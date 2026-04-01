#!/bin/bash
mkdir -p /tmp/compose-health
cat > /tmp/compose-health/docker-compose.yaml << 'EOF'
version: "3.8"
services:
  db:
    image: redis:7-alpine
    container_name: health-db
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5
  web:
    image: nginx:alpine
    container_name: health-web
    depends_on:
      db:
        condition: service_healthy
EOF
cd /tmp/compose-health && docker compose up -d
