#!/bin/bash
for f in /tmp/custom_ua.txt /tmp/cookies.txt /tmp/verbose.txt; do
    if [ ! -f "$f" ]; then
        echo "FAIL: $f not found"
        exit 1
    fi
done
if [ -s /tmp/verbose.txt ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: verbose output is empty"
exit 1
