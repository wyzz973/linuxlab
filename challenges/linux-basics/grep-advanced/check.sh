#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
# Should have line numbers and contain @ symbols
if grep -q "@" /tmp/result.txt && grep -qE "^[0-9]+:" /tmp/result.txt; then
    count=$(wc -l < /tmp/result.txt | tr -d ' ')
    if [ "$count" -ge 3 ]; then
        echo "PASS"
        exit 0
    fi
fi
echo "FAIL: Expected lines with line numbers containing email addresses"
exit 1
