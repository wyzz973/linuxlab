#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "sleep" /tmp/result.txt; then
    echo "PASS"
    exit 0
else
    echo "FAIL: No sleep job found in output"
    exit 1
fi
