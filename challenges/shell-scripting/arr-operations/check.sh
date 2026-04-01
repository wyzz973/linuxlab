#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "Students: 10\nHighest: 97\nLowest: 70\nTotal: 853\nAverage: 85")
if [ "$output" = "$expected" ]; then
    echo "PASS: 成绩统计正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
