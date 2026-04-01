#!/bin/bash
mkdir -p /tmp/compose-vol/config
cat > /tmp/compose-vol/docker-compose.yaml << 'EOF'
version: "3.8"
services:
  web:
    image: nginx:alpine
    container_name: vol-web
    ports:
      - "8094:80"
    volumes:
      - web-data:/usr/share/nginx/html
      - ./config:/etc/nginx/conf.d
volumes:
  web-data:
EOF
cd /tmp/compose-vol && docker compose up -d
