#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
count=$(echo "$output" | grep -c '^\(Today:\|User:\|Lines in passwd:\)')
if [ "$count" -eq 3 ]; then
    echo "PASS: 命令替换使用正确"
    exit 0
else
    echo "FAIL: 期望输出3行，分别以 Today:, User:, Lines in passwd: 开头"
    echo "实际输出: '$output'"
    exit 1
fi
