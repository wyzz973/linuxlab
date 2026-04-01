#!/bin/bash
if ! docker ps --format '{{.Names}}' | grep -q '^web-server$'; then
    echo "FAIL: 容器 web-server 未在运行"
    exit 1
fi
PORT_MAP=$(docker port web-server 80 2>/dev/null)
if echo "$PORT_MAP" | grep -q "8080"; then
    echo "PASS"
    exit 0
fi
echo "FAIL: 端口映射不正确，期望 80->8080"
exit 1
