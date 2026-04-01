#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "apple\ncherry\napple banana cherry date elderberry\napple blueberry cherry date elderberry")
if [ "$output" = "$expected" ]; then
    echo "PASS: 数组定义与访问正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
