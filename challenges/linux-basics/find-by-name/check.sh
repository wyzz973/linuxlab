#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
count=$(grep -c "\.py$" /tmp/result.txt)
if [ "$count" -ge 3 ]; then
    echo "PASS: Found $count Python files"
    exit 0
else
    echo "FAIL: Expected at least 3 .py files, found $count"
    exit 1
fi
