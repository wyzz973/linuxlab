#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "is file\nnot directory\nis readable\nhas content\nis directory\nis executable\nnot exists")
if [ "$output" = "$expected" ]; then
    echo "PASS: 文件测试运算符使用正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
