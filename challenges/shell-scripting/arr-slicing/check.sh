#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "30 40 50\n11\n20 30 40 50 60 70 80 90 100 110")
if [ "$output" = "$expected" ]; then
    echo "PASS: 数组切片正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
