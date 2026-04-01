#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "hello: 5\nworld: 5\nshell: 5")
if [ "$output" = "$expected" ]; then
    echo "PASS: 字符串长度获取正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
