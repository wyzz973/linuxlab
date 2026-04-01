#!/bin/bash
mkdir -p /tmp/multi-app
cat > /tmp/multi-app/docker-compose.yaml << 'EOF'
version: "3.8"
services:
  web:
    image: nginx:alpine
    ports:
      - "8093:80"
  api:
    image: node:18-alpine
    command: node -e "require('http').createServer((req,res)=>{res.end('API OK')}).listen(3000)"
    ports:
      - "3000:3000"
  db:
    image: redis:7-alpine
    ports:
      - "6379:6379"
EOF
docker compose -f /tmp/multi-app/docker-compose.yaml up -d
