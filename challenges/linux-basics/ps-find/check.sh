#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "PID" /tmp/result.txt && grep -q "COMMAND\|CMD" /tmp/result.txt; then
    lines=$(wc -l < /tmp/result.txt | tr -d ' ')
    if [ "$lines" -gt 2 ]; then
        echo "PASS"
        exit 0
    fi
fi
echo "FAIL: Output doesn't look like ps aux output"
exit 1
