#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
line1=$(echo "$output" | sed -n '1p')
line3=$(echo "$output" | sed -n '3p')
if [ "$line1" = "Hello from function!" ] && [ "$line3" = "30" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
