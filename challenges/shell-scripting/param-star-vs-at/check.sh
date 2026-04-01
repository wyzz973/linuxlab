#!/bin/bash
output=$(bash /home/learner/solution.sh "hello world" "foo bar" "test" 2>/dev/null)
line1=$(echo "$output" | sed -n '1p')
line2=$(echo "$output" | sed -n '2p')
if echo "$line1" | grep -q "1 iterations" && echo "$line2" | grep -q "3 iterations"; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
