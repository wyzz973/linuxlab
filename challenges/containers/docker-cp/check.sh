#!/bin/bash
if ! docker ps --format '{{.Names}}' | grep -q '^cp-demo$'; then
    echo "FAIL: 容器 cp-demo 未在运行"
    exit 1
fi
# Check file was copied to container
CONTENT=$(docker exec cp-demo cat /tmp/upload.txt 2>/dev/null)
if ! echo "$CONTENT" | grep -q "file from host"; then
    echo "FAIL: 文件未正确复制到容器"
    exit 1
fi
# Check file was copied from container
if [ ! -f /tmp/nginx.conf ]; then
    echo "FAIL: /tmp/nginx.conf 不存在（未从容器中复制）"
    exit 1
fi
if grep -q "nginx" /tmp/nginx.conf; then
    echo "PASS"
    exit 0
fi
echo "FAIL: /tmp/nginx.conf 内容不正确"
exit 1
