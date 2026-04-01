#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "1 2 3 4 5\n2\n4\n6\n8\n10\n5050")
if [ "$output" = "$expected" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
