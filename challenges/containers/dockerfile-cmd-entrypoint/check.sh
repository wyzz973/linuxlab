#!/bin/bash
if [ ! -f /tmp/greeter/Dockerfile ]; then
    echo "FAIL: /tmp/greeter/Dockerfile 不存在"
    exit 1
fi
if ! docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "greeter:v1"; then
    echo "FAIL: 镜像 greeter:v1 未构建"
    exit 1
fi
OUTPUT1=$(docker run --rm greeter:v1 2>/dev/null)
if ! echo "$OUTPUT1" | grep -q "Hello World"; then
    echo "FAIL: docker run greeter:v1 输出应为 'Hello World'，实际为 '$OUTPUT1'"
    exit 1
fi
OUTPUT2=$(docker run --rm greeter:v1 Docker 2>/dev/null)
if ! echo "$OUTPUT2" | grep -q "Hello Docker"; then
    echo "FAIL: docker run greeter:v1 Docker 输出应为 'Hello Docker'，实际为 '$OUTPUT2'"
    exit 1
fi
echo "PASS"
exit 0
