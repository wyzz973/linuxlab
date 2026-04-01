#!/bin/bash
for f in /tmp/recent_logs.txt /tmp/kernel_logs.txt /tmp/units.txt; do
    if [ ! -f "$f" ]; then
        echo "FAIL: $f not found"
        exit 1
    fi
done
# At least recent_logs should have content
if [ -s /tmp/recent_logs.txt ] || [ -s /tmp/kernel_logs.txt ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: log files are empty"
exit 1
