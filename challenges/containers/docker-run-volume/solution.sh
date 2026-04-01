#!/bin/bash
mkdir -p /tmp/web-content
echo "Hello Docker" > /tmp/web-content/index.html
docker run -d --name web-mount -p 8081:80 -v /tmp/web-content:/usr/share/nginx/html nginx
