#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
count=$(echo "$output" | grep -c '\(Script name:\|PID:\|Exit status:\)')
if [ "$count" -ge 4 ]; then
    echo "PASS: 特殊变量使用正确"
    exit 0
else
    echo "FAIL: 期望包含 Script name:, PID:, Exit status: 共4行"
    echo "实际输出: '$output'"
    exit 1
fi
