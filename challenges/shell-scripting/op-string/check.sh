#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "not equal\ndifferent\nstr3 is empty\nstr1 is not empty\nstr1 has value")
if [ "$output" = "$expected" ]; then
    echo "PASS: 字符串运算符使用正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
