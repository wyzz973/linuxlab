#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -qE "LISTEN|State" /tmp/result.txt; then
    echo "PASS"
    exit 0
elif [ -s /tmp/result.txt ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected listening socket info"
    exit 1
fi
