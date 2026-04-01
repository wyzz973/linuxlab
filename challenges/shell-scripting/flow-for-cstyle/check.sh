#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "1^2 = 1\n2^2 = 4\n3^2 = 9\n4^2 = 16\n5^2 = 25\n5x1=5 5x2=10 5x3=15 5x4=20 5x5=25")
if [ "$output" = "$expected" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
