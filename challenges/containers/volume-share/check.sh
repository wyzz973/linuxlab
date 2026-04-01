#!/bin/bash
if ! docker volume ls --format '{{.Name}}' | grep -q '^shared-data$'; then
    echo "FAIL: 数据卷 shared-data 不存在"
    exit 1
fi
if ! docker ps --format '{{.Names}}' | grep -q '^producer$'; then
    echo "FAIL: producer 容器未在运行"
    exit 1
fi
if ! docker ps --format '{{.Names}}' | grep -q '^consumer$'; then
    echo "FAIL: consumer 容器未在运行"
    exit 1
fi
# Verify shared volume works
CONTENT=$(docker exec consumer cat /shared/timestamp.txt 2>/dev/null)
if [ -z "$CONTENT" ]; then
    echo "FAIL: consumer 无法读取共享数据"
    exit 1
fi
if [ ! -f /tmp/shared-data.txt ]; then
    echo "FAIL: /tmp/shared-data.txt 不存在"
    exit 1
fi
echo "PASS"
exit 0
