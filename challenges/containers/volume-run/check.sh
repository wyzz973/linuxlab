#!/bin/bash
if ! docker volume ls --format '{{.Name}}' | grep -q '^persist-data$'; then
    echo "FAIL: 数据卷 persist-data 不存在"
    exit 1
fi
# writer should be removed
if docker ps -a --format '{{.Names}}' | grep -q '^writer$'; then
    echo "FAIL: writer 容器应该已被删除"
    exit 1
fi
if [ ! -f /tmp/persist-check.txt ]; then
    echo "FAIL: /tmp/persist-check.txt 不存在"
    exit 1
fi
if grep -q "data survives" /tmp/persist-check.txt; then
    echo "PASS"
    exit 0
fi
echo "FAIL: 数据未正确持久化"
exit 1
