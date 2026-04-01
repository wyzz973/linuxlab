#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "5\n3\n17\n42\n61\n85\n85\n3")
if [ "$output" = "$expected" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
