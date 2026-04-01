#!/bin/bash
if ! docker ps -a --format '{{.Names}}' | grep -q '^interactive-ubuntu$'; then
    echo "FAIL: 容器 interactive-ubuntu 不存在"
    exit 1
fi
CONTENT=$(docker cp interactive-ubuntu:/data/hello.txt /tmp/check-hello.txt 2>&1)
if [ $? -ne 0 ]; then
    echo "FAIL: 容器中 /data/hello.txt 文件不存在"
    exit 1
fi
if grep -q "interactive mode" /tmp/check-hello.txt; then
    echo "PASS"
    rm -f /tmp/check-hello.txt
    exit 0
fi
echo "FAIL: 文件内容不正确"
rm -f /tmp/check-hello.txt
exit 1
