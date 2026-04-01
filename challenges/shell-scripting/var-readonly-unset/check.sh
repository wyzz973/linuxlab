#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "runoob.com\n")
if [ "$output" = "$expected" ]; then
    echo "PASS: readonly 和 unset 使用正确"
    exit 0
else
    echo "FAIL: 期望输出 'runoob.com' 后跟空行，实际输出 '$output'"
    exit 1
fi
