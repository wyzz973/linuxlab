#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "1 2 4 5 7 8 10 11 13 14 16 17 19 20\nStopped at 20, sum=210")
if [ "$output" = "$expected" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
