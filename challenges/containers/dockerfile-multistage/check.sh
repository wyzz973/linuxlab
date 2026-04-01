#!/bin/bash
if ! docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "multistage:v1"; then
    echo "FAIL: 镜像 multistage:v1 未构建"
    exit 1
fi
# Check image size is small (should be based on alpine, not gcc)
SIZE=$(docker inspect multistage:v1 --format '{{.Size}}')
# gcc image is ~1GB, alpine-based should be much smaller
if [ "$SIZE" -gt 500000000 ]; then
    echo "FAIL: 镜像太大 (${SIZE} bytes)，可能未使用多阶段构建"
    exit 1
fi
OUTPUT=$(docker run --rm multistage:v1 2>/dev/null)
if echo "$OUTPUT" | grep -q "Hello from multi-stage"; then
    echo "PASS"
    exit 0
fi
echo "FAIL: 运行输出不正确"
exit 1
