#!/bin/bash
for name in prune-test-1 prune-test-2 prune-test-3; do
    if docker ps -a --format '{{.Names}}' | grep -q "^${name}$"; then
        echo "FAIL: 容器 ${name} 仍然存在"
        exit 1
    fi
done
echo "PASS"
exit 0
