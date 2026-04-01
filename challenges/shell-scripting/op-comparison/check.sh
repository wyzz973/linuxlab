#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "a == b: false\na != b: true\na > b: false\na < b: true\na >= 10: true\nb <= 20: true")
if [ "$output" = "$expected" ]; then
    echo "PASS: 比较运算符使用正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
