#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
first=$(head -1 /tmp/result.txt | awk '{print $1}')
if [ "$first" = "Alice" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: First line should be Alice (highest math score 98), got $first"
    exit 1
fi
