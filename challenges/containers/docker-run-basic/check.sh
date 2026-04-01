#!/bin/bash
if [ -f /tmp/hello-docker.txt ]; then
    if grep -q "Hello from Docker" /tmp/hello-docker.txt; then
        echo "PASS"
        exit 0
    fi
    echo "FAIL: /tmp/hello-docker.txt 中没有包含预期的 hello-world 输出"
    exit 1
fi
echo "FAIL: /tmp/hello-docker.txt 文件不存在"
exit 1
