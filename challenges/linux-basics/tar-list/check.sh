#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "secret.txt" /tmp/result.txt; then
    echo "PASS"
    exit 0
else
    echo "FAIL: secret.txt not found in listing"
    exit 1
fi
