#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "localhost" /tmp/result.txt; then
    echo "FAIL: 'localhost' still found in output"
    exit 1
fi
localhost_count=$(grep -c "db.production.com" /tmp/result.txt)
if [ "$localhost_count" -ge 3 ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected at least 3 replacements"
    exit 1
fi
