#!/bin/bash
if [ -f /tmp/page.html ] && [ -s /tmp/page.html ]; then
    echo "PASS"
    exit 0
fi
if [ -f /tmp/result.txt ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: No output file found"
exit 1
