#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "Filesystem" /tmp/result.txt && grep -q "Type" /tmp/result.txt; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Output doesn't contain expected df -hT headers"
    exit 1
fi
