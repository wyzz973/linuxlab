#!/bin/bash
if ! docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "redis:7-alpine"; then
    echo "FAIL: redis:7-alpine 镜像未拉取"
    exit 1
fi
if ! docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "nginx:1.25"; then
    echo "FAIL: nginx:1.25 镜像未拉取"
    exit 1
fi
if [ ! -f /tmp/image-list.txt ]; then
    echo "FAIL: /tmp/image-list.txt 不存在"
    exit 1
fi
echo "PASS"
exit 0
