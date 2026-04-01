#!/bin/bash
if ! docker ps --format '{{.Names}}' | grep -q '^log-demo$'; then
    echo "FAIL: 容器 log-demo 未在运行"
    exit 1
fi
if [ ! -f /tmp/container-logs.txt ]; then
    echo "FAIL: /tmp/container-logs.txt 不存在"
    exit 1
fi
if [ ! -f /tmp/container-ip.txt ]; then
    echo "FAIL: /tmp/container-ip.txt 不存在"
    exit 1
fi
if [ ! -s /tmp/container-ip.txt ]; then
    echo "FAIL: /tmp/container-ip.txt 为空"
    exit 1
fi
echo "PASS"
exit 0
