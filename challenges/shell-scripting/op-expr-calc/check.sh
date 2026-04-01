#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "10 + 5 = 15\n20 - 8 = 12\n6 * 7 = 42\n100 / 4 = 25\n17 %% 3 = 2")
if [ "$output" = "$expected" ]; then
    echo "PASS: 计算器正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
