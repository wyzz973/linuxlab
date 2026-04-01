#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "home/user/documents/report.tar.gz\nreport.tar.gz\ngz")
if [ "$output" = "$expected" ]; then
    echo "PASS: 左侧模式匹配正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
