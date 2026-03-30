#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
lines=$(wc -l < /tmp/result.txt | tr -d ' ')
if [ "$lines" -eq 11 ] && grep -q "Line 20" /tmp/result.txt && grep -q "Line 30" /tmp/result.txt; then
    if ! grep -q "Line 19:" /tmp/result.txt && ! grep -q "Line 31:" /tmp/result.txt; then
        echo "PASS"
        exit 0
    fi
fi
echo "FAIL: Expected lines 20-30 (11 lines)"
exit 1
