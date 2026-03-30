#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "main.py" /tmp/result.txt && grep -q "test.py" /tmp/result.txt && grep -q "guide.md" /tmp/result.txt; then
    echo "PASS"
    exit 0
fi
echo "FAIL: Expected recursive listing"
exit 1
