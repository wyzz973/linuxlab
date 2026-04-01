#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
line2=$(echo "$output" | sed -n '2p')
line3=$(echo "$output" | sed -n '3p')
if [ "$line2" = "LinuxLab" ] && [ "$line3" = "PATH contains /usr/bin" ]; then
    echo "PASS: 环境变量使用正确"
    exit 0
else
    echo "FAIL: 输出不符合预期"
    echo "实际输出: '$output'"
    exit 1
fi
