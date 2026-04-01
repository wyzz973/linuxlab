#!/bin/bash
if [ ! -f /tmp/ignore-demo/.dockerignore ]; then
    echo "FAIL: .dockerignore 文件不存在"
    exit 1
fi
if ! docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "ignore-demo:v1"; then
    echo "FAIL: 镜像 ignore-demo:v1 未构建"
    exit 1
fi
# Check excluded files are not in image
if docker run --rm ignore-demo:v1 ls /app/test.sh 2>/dev/null; then
    echo "FAIL: test.sh 不应该在镜像中"
    exit 1
fi
if docker run --rm ignore-demo:v1 ls /app/README.md 2>/dev/null; then
    echo "FAIL: README.md 不应该在镜像中"
    exit 1
fi
# Check app.sh is in image
if docker run --rm ignore-demo:v1 ls /app/app.sh 2>/dev/null; then
    echo "PASS"
    exit 0
fi
echo "FAIL: app.sh 应该在镜像中"
exit 1
