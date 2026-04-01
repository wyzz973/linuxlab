#!/bin/bash
if ! docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "workdir-app:v1"; then
    echo "FAIL: 镜像 workdir-app:v1 未构建"
    exit 1
fi
WORKDIR=$(docker inspect workdir-app:v1 --format '{{.Config.WorkingDir}}')
if [ "$WORKDIR" != "/app" ]; then
    echo "FAIL: 工作目录不是 /app，当前为 $WORKDIR"
    exit 1
fi
OUTPUT=$(docker run --rm workdir-app:v1 2>/dev/null)
if echo "$OUTPUT" | grep -q "app is running"; then
    echo "PASS"
    exit 0
fi
echo "FAIL: 运行输出不正确，期望 'app is running'，实际为 '$OUTPUT'"
exit 1
