#!/bin/bash
if ! docker ps --format '{{.Names}}' | grep -q 'updown-web'; then
    echo "FAIL: updown-web 容器未在运行"
    exit 1
fi
if ! docker ps --format '{{.Names}}' | grep -q 'updown-cache'; then
    echo "FAIL: updown-cache 容器未在运行"
    exit 1
fi
if [ ! -f /tmp/compose-status.txt ]; then
    echo "FAIL: /tmp/compose-status.txt 不存在"
    exit 1
fi
echo "PASS"
exit 0
