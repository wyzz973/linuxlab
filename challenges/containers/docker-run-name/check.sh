#!/bin/bash
if [ ! -f /tmp/container-test.txt ]; then
    echo "FAIL: /tmp/container-test.txt 不存在"
    exit 1
fi
if ! grep -q "container test" /tmp/container-test.txt; then
    echo "FAIL: 文件内容不正确"
    exit 1
fi
if docker ps -a --format '{{.Names}}' | grep -q '^temp-container$'; then
    echo "FAIL: 容器未被自动删除（--rm 参数未生效）"
    exit 1
fi
echo "PASS"
exit 0
