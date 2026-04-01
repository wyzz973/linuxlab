#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "13\n7\n30\n3\n1\n26")
if [ "$output" = "$expected" ]; then
    echo "PASS: 算术运算正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
