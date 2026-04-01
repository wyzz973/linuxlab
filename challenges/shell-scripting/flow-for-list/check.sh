#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
expected=$(printf "Linux\nmacOS\nWindows\nAndroid\niOS\n1. apple\n2. banana\n3. cherry")
if [ "$output" = "$expected" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
