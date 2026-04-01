#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
lines=$(echo "$output" | wc -l | tr -d ' ')
line1=$(echo "$output" | sed -n '1p')
line9=$(echo "$output" | sed -n '9p')
if [ "$lines" = "9" ] && echo "$line1" | grep -q "1x1=1" && echo "$line9" | grep -q "9x9=81"; then
    echo "PASS"
    exit 0
else
    echo "FAIL: 期望9行乘法表"
    echo "实际行数: $lines"
    exit 1
fi
