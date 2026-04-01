#!/bin/bash
if ! docker ps --format '{{.Names}}' | grep -q 'vol-web'; then
    echo "FAIL: vol-web 容器未在运行"
    exit 1
fi
MOUNTS=$(docker inspect vol-web --format '{{range .Mounts}}{{.Type}}:{{.Destination}} {{end}}')
if ! echo "$MOUNTS" | grep -q "/usr/share/nginx/html"; then
    echo "FAIL: web-data 卷未挂载到 /usr/share/nginx/html"
    exit 1
fi
echo "PASS"
exit 0
