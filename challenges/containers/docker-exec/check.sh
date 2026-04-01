#!/bin/bash
if ! docker ps --format '{{.Names}}' | grep -q '^exec-demo$'; then
    echo "FAIL: 容器 exec-demo 未在运行"
    exit 1
fi
# Check file inside container
CONTENT=$(docker exec exec-demo cat /tmp/exec-test.txt 2>/dev/null)
if ! echo "$CONTENT" | grep -q "executed inside container"; then
    echo "FAIL: 容器内 /tmp/exec-test.txt 文件不存在或内容不正确"
    exit 1
fi
if [ ! -f /tmp/container-processes.txt ]; then
    echo "FAIL: /tmp/container-processes.txt 不存在"
    exit 1
fi
echo "PASS"
exit 0
