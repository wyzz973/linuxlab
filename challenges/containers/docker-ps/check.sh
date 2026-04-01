#!/bin/bash
if ! docker ps --format '{{.Names}}' | grep -q '^running-web$'; then
    echo "FAIL: 容器 running-web 未在运行"
    exit 1
fi
if ! docker ps -a --format '{{.Names}}' | grep -q '^stopped-task$'; then
    echo "FAIL: 容器 stopped-task 不存在"
    exit 1
fi
if [ ! -f /tmp/all-containers.txt ]; then
    echo "FAIL: /tmp/all-containers.txt 不存在"
    exit 1
fi
if [ ! -f /tmp/running-containers.txt ]; then
    echo "FAIL: /tmp/running-containers.txt 不存在"
    exit 1
fi
echo "PASS"
exit 0
