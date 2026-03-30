#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
# Check that line numbers are present
if grep -qE "^\s+1\s" /tmp/result.txt && grep -qE "^\s+[2-8]\s" /tmp/result.txt; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Line numbers not found in output"
    exit 1
fi
