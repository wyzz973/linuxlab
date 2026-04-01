#!/bin/bash
output=$(bash /home/learner/solution.sh 2>/dev/null)
has_hello=$(echo "$output" | grep -c "Hello, World")
has_price=$(echo "$output" | grep -c 'Price: \$100')
has_server=$(echo "$output" | grep -c '\[server\]')
has_port=$(echo "$output" | grep -c 'port=8080')
if [ "$has_hello" -ge 1 ] && [ "$has_price" -ge 1 ] && [ "$has_server" -ge 1 ] && [ "$has_port" -ge 1 ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    echo "实际输出: '$output'"
    exit 1
fi
