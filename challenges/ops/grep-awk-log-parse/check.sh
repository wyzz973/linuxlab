#!/bin/bash
if [ ! -f /tmp/errors.txt ] || [ ! -s /tmp/errors.txt ]; then
    echo "FAIL: /tmp/errors.txt not found or empty"
    exit 1
fi
ERROR_COUNT=$(wc -l < /tmp/errors.txt | tr -d " ")
if [ "$ERROR_COUNT" -ne 7 ]; then
    echo "FAIL: expected 7 error lines, got $ERROR_COUNT"
    exit 1
fi
if [ ! -f /tmp/error_summary.txt ] || [ ! -s /tmp/error_summary.txt ]; then
    echo "FAIL: /tmp/error_summary.txt not found or empty"
    exit 1
fi
if [ ! -f /tmp/hourly_errors.txt ] || [ ! -s /tmp/hourly_errors.txt ]; then
    echo "FAIL: /tmp/hourly_errors.txt not found or empty"
    exit 1
fi
echo "PASS"
exit 0
