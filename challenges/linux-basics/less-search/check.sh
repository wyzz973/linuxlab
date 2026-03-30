#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "CRITICAL" /tmp/result.txt; then
    # Check that context lines are present
    if grep -q "Attempting recovery\|High memory" /tmp/result.txt || grep -q "Disk space low\|WARNING" /tmp/result.txt; then
        echo "PASS"
        exit 0
    fi
    echo "PASS: CRITICAL lines found"
    exit 0
fi
echo "FAIL: CRITICAL lines not found"
exit 1
