#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "Adult with passing score\nNot (under 18 or above 90)\nage is not 30\nWorking age\nNot (excellent or young)")
if [ "$output" = "$expected" ]; then
    echo "PASS: 布尔运算符使用正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
