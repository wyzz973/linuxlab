#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "app.log" /tmp/result.txt && grep -q "access.log" /tmp/result.txt; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Recent .log files not found"
    exit 1
fi
