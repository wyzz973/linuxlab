#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
line1=$(echo "$output" | sed -n '1p')
line2=$(echo "$output" | sed -n '2p')
line3=$(echo "$output" | sed -n '3p')
if [ "$line1" = "I am local" ] && [ "$line2" = "Modified in function" ] && [ -z "$line3" ]; then
    echo "PASS: 变量作用域理解正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
