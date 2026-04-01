#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "Inside: count=1, local_count=100\nInside: count=2, local_count=100\nInside: count=3, local_count=100\nOutside: count=3\nOutside: local_count=")
if [ "$output" = "$expected" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
