#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
count=$(grep -c "TODO" /tmp/result.txt)
if [ "$count" -ge 4 ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected at least 4 TODO matches, found $count"
    exit 1
fi
