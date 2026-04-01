#!/bin/bash
mkdir -p /tmp/compose-env
cat > /tmp/compose-env/.env << 'EOF'
MYSQL_ROOT_PASSWORD=rootpass123
MYSQL_DATABASE=mydb
EOF
cat > /tmp/compose-env/docker-compose.yaml << 'EOF'
version: "3.8"
services:
  db:
    image: mysql:8.0
    container_name: compose-mysql
    environment:
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
      - MYSQL_DATABASE=${MYSQL_DATABASE}
EOF
cd /tmp/compose-env && docker compose up -d
