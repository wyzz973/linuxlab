#!/bin/bash
if ! docker volume ls --format '{{.Name}}' | grep -q '^named-vol$'; then
    echo "FAIL: 命名卷 named-vol 不存在"
    exit 1
fi
# Check named volume has data
CONTENT=$(docker run --rm -v named-vol:/data alpine cat /data/test.txt 2>/dev/null)
if ! echo "$CONTENT" | grep -q "named volume data"; then
    echo "FAIL: 命名卷中的数据不正确"
    exit 1
fi
# Check bind mount
if [ ! -f /tmp/bind-data/test.txt ]; then
    echo "FAIL: /tmp/bind-data/test.txt 不存在"
    exit 1
fi
if ! grep -q "bind mount data" /tmp/bind-data/test.txt; then
    echo "FAIL: 绑定挂载中的数据不正确"
    exit 1
fi
echo "PASS"
exit 0
