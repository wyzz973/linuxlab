#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
lines=$(wc -l < /tmp/result.txt | tr -d ' ')
if [ "$lines" -eq 5 ]; then
    first=$(head -1 /tmp/result.txt | tr -s ' ')
    if echo "$first" | grep -q "192.168.1.100"; then
        echo "PASS"
        exit 0
    fi
fi
echo "FAIL: Expected 5 lines with 192.168.1.100 first"
exit 1
