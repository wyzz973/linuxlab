#!/bin/bash
if [ ! -f /tmp/webapp/Dockerfile ]; then
    echo "FAIL: /tmp/webapp/Dockerfile 不存在"
    exit 1
fi
if ! grep -q "COPY" /tmp/webapp/Dockerfile; then
    echo "FAIL: Dockerfile 中未使用 COPY 指令"
    exit 1
fi
if ! docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "webapp:v1"; then
    echo "FAIL: 镜像 webapp:v1 未构建"
    exit 1
fi
# Verify the content is correct
CONTENT=$(docker run --rm webapp:v1 cat /usr/share/nginx/html/index.html 2>/dev/null)
if echo "$CONTENT" | grep -q "Docker Web App"; then
    echo "PASS"
    exit 0
fi
echo "FAIL: 镜像中 index.html 内容不正确"
exit 1
