#!/bin/bash
if [ ! -f /tmp/myapp/Dockerfile ]; then
    echo "FAIL: /tmp/myapp/Dockerfile 不存在"
    exit 1
fi
if ! grep -q "FROM ubuntu:22.04" /tmp/myapp/Dockerfile; then
    echo "FAIL: Dockerfile 中未使用 ubuntu:22.04 作为基础镜像"
    exit 1
fi
if ! grep -q "curl" /tmp/myapp/Dockerfile; then
    echo "FAIL: Dockerfile 中未安装 curl"
    exit 1
fi
if docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "myapp:v1"; then
    echo "PASS"
    exit 0
fi
echo "FAIL: 镜像 myapp:v1 未构建"
exit 1
