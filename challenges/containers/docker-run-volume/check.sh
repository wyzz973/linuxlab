#!/bin/bash
if [ ! -f /tmp/web-content/index.html ]; then
    echo "FAIL: /tmp/web-content/index.html 不存在"
    exit 1
fi
if ! grep -q "Hello Docker" /tmp/web-content/index.html; then
    echo "FAIL: index.html 内容不正确"
    exit 1
fi
if ! docker ps --format '{{.Names}}' | grep -q '^web-mount$'; then
    echo "FAIL: 容器 web-mount 未在运行"
    exit 1
fi
MOUNTS=$(docker inspect web-mount --format '{{range .Mounts}}{{.Source}}:{{.Destination}} {{end}}')
if echo "$MOUNTS" | grep -q "/tmp/web-content:/usr/share/nginx/html"; then
    echo "PASS"
    exit 0
fi
echo "FAIL: 目录挂载不正确"
exit 1
