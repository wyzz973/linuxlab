#!/bin/bash
if ! docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "labeled-app:v2"; then
    echo "FAIL: 镜像 labeled-app:v2 未构建"
    exit 1
fi
LABELS=$(docker inspect labeled-app:v2 --format '{{json .Config.Labels}}')
if ! echo "$LABELS" | grep -q "maintainer"; then
    echo "FAIL: 镜像缺少 maintainer 标签"
    exit 1
fi
if ! echo "$LABELS" | grep -q "version"; then
    echo "FAIL: 镜像缺少 version 标签"
    exit 1
fi
echo "PASS"
exit 0
