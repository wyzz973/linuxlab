#!/bin/bash
output=$(bash /home/learner/solution.sh apple banana cherry date 2>/dev/null)
expected=$(printf "Count: 4\napple\nbanana\ncherry\ndate\nAll: apple,banana,cherry,date")
if [ "$output" = "$expected" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
