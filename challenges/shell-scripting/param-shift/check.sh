#!/bin/bash
output=$(bash /home/learner/solution.sh one two three four 2>/dev/null)
expected=$(printf "Processing #1: one (3 remaining)\nProcessing #2: two (2 remaining)\nProcessing #3: three (1 remaining)\nProcessing #4: four (0 remaining)")
if [ "$output" = "$expected" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
