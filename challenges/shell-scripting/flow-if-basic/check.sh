#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
if [ "$output" = "42 is positive" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: 期望 '42 is positive'，实际 '$output'"
    exit 1
fi
