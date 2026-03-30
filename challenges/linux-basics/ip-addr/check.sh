#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "lo\|inet\|127.0.0.1" /tmp/result.txt; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected network interface info"
    exit 1
fi
