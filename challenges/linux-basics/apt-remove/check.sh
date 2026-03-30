#!/bin/bash
if [ -f /tmp/result.txt ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
