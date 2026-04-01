#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
if [ "$output" = "Welcome to LinuxLab v3" ]; then
    echo "PASS: 变量定义和使用正确"
    exit 0
else
    echo "FAIL: 期望输出 'Welcome to LinuxLab v3'，实际输出 '$output'"
    exit 1
fi
