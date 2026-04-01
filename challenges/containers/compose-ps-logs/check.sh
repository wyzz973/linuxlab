#!/bin/bash
if ! docker ps --format '{{.Names}}' | grep -q 'logs-web'; then
    echo "FAIL: logs-web 容器未在运行"
    exit 1
fi
if [ ! -f /tmp/compose-web-logs.txt ]; then
    echo "FAIL: /tmp/compose-web-logs.txt 不存在"
    exit 1
fi
if [ ! -f /tmp/compose-ps.txt ]; then
    echo "FAIL: /tmp/compose-ps.txt 不存在"
    exit 1
fi
echo "PASS"
exit 0
