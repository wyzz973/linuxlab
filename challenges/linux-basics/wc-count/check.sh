#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
count=$(tr -d ' \n' < /tmp/result.txt)
if [ "$count" = "20" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected 20, got '$count'"
    exit 1
fi
