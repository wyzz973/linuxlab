#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -qi "mem\|Mem\|total" /tmp/result.txt; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected memory info"
    exit 1
fi
