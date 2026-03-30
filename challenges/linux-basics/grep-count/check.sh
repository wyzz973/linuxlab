#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
count=$(tr -d ' \n' < /tmp/result.txt)
if [ "$count" = "5" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected 5, got '$count'"
    exit 1
fi
