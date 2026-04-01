#!/bin/bash
if [ ! -f /tmp/os-release.txt ]; then
    echo "FAIL: /tmp/os-release.txt 不存在"
    exit 1
fi
if grep -q "Alpine" /tmp/os-release.txt; then
    echo "PASS"
    exit 0
fi
echo "FAIL: 文件内容不正确，应该包含 Alpine 的系统信息"
exit 1
