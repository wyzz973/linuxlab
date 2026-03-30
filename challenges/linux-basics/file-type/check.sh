#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "text" /tmp/result.txt; then
    echo "PASS"
    exit 0
else
    echo "FAIL: File type information not found"
    exit 1
fi
