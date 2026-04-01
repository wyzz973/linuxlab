#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "I love banana, apple is sweet\nI love orange, orange is sweet\n-home-user-documents-file.txt")
if [ "$output" = "$expected" ]; then
    echo "PASS: 字符串替换正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
