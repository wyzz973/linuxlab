#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "Alice is 28 years old, lives in Beijing\nBob is 35 years old, lives in Shanghai\nCharlie is 22 years old, lives in Guangzhou")
if [ "$output" = "$expected" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
