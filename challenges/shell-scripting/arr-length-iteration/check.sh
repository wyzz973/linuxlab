#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
line1=$(echo "$output" | sed -n '1p')
line7=$(echo "$output" | sed -n '7p')
line11=$(echo "$output" | sed -n '11p')
total=$(echo "$output" | wc -l | tr -d ' ')
if [ "$line1" = "5" ] && [ "$line7" = "0: red" ] && [ "$line11" = "4: purple" ] && [ "$total" = "11" ]; then
    echo "PASS: 数组遍历正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
