#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "Beijing\nChina\nFrance\nJapan\nBeijing\nParis\nTokyo\n3")
if [ "$output" = "$expected" ]; then
    echo "PASS: 关联数组使用正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
