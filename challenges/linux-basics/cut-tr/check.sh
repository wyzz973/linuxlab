#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
# Check lowercase conversion
if grep -q "Root\|Admin\|Guest\|MySQL\|Daemon" /tmp/result.txt; then
    echo "FAIL: Uppercase letters still present"
    exit 1
fi
if grep -q "root:/bin/bash" /tmp/result.txt && grep -q "admin:/bin/bash" /tmp/result.txt; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected lowercase username:shell format"
    exit 1
fi
