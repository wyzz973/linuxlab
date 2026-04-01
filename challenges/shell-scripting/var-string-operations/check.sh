#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
line1=$(echo "$output" | sed -n '1p')
line2=$(echo "$output" | sed -n '2p')
if [ "$line1" = "46" ] && [ "$line2" = "runoob" ]; then
    echo "PASS: 字符串操作正确"
    exit 0
else
    echo "FAIL: 期望第一行 '46'，第二行 'runoob'"
    echo "实际输出: '$output'"
    exit 1
fi
