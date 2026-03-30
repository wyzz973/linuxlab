#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "4 packets transmitted" /tmp/result.txt; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected 4 packets transmitted"
    exit 1
fi
