#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
line1=$(echo "$output" | sed -n '1p')
line2=$(echo "$output" | sed -n '2p')
if [ "$line1" = 'Hello, Shell!' ] && [ "$line2" = 'Hello, $lang!' ]; then
    echo "PASS: 引号规则理解正确"
    exit 0
else
    echo "FAIL: 期望第一行 'Hello, Shell!' 第二行 'Hello, \$lang!'"
    echo "实际输出: '$output'"
    exit 1
fi
