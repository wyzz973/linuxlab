#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
line1=$(echo "$output" | sed -n '1p')
last=$(echo "$output" | tail -1)
if [ "$line1" = "Sum: 55, stopped at: 10" ] && [ "$last" = "2" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
