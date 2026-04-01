#!/bin/bash
if [ ! -f /tmp/bridge-info.txt ]; then
    echo "FAIL: /tmp/bridge-info.txt 不存在"
    exit 1
fi
if ! docker ps --format '{{.Names}}' | grep -q 'bridge-test-1'; then
    echo "FAIL: bridge-test-1 容器未在运行"
    exit 1
fi
if ! docker ps --format '{{.Names}}' | grep -q 'bridge-test-2'; then
    echo "FAIL: bridge-test-2 容器未在运行"
    exit 1
fi
if [ ! -f /tmp/container-ips.txt ]; then
    echo "FAIL: /tmp/container-ips.txt 不存在"
    exit 1
fi
echo "PASS"
exit 0
