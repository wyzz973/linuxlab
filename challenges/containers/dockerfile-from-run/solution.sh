#!/bin/bash
mkdir -p /tmp/myapp
cat > /tmp/myapp/Dockerfile << 'EOF'
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*
EOF
docker build -t myapp:v1 /tmp/myapp/
