#!/bin/bash
if [ ! -f /tmp/date_result.txt ]; then
    echo "FAIL: /tmp/date_result.txt not found"
    exit 1
fi
# Check date format YYYY-MM-DD HH:MM:SS
if grep -qE "^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}" /tmp/date_result.txt; then
    if [ -f /tmp/cal_result.txt ]; then
        echo "PASS"
        exit 0
    else
        echo "FAIL: /tmp/cal_result.txt not found"
        exit 1
    fi
fi
echo "FAIL: Date format should be YYYY-MM-DD HH:MM:SS"
exit 1
