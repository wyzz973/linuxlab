#!/bin/bash
output=$(bash /home/learner/solution.sh Alice 25 Beijing 2>/dev/null)
expected=$(printf "Name: Alice\nAge: 25\nCity: Beijing\nScript: /home/learner/solution.sh")
if [ "$output" = "$expected" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
