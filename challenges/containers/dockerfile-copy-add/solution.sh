#!/bin/bash
mkdir -p /tmp/webapp
echo '<h1>Docker Web App</h1>' > /tmp/webapp/index.html
cat > /tmp/webapp/Dockerfile << 'EOF'
FROM nginx:alpine
COPY index.html /usr/share/nginx/html/
EOF
docker build -t webapp:v1 /tmp/webapp/
