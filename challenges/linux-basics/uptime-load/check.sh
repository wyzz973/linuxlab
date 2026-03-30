#!/bin/bash
if [ ! -f /tmp/uptime_result.txt ]; then
    echo "FAIL: /tmp/uptime_result.txt not found"
    exit 1
fi
if grep -qi "load average\|up" /tmp/uptime_result.txt; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected uptime info"
    exit 1
fi
