#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "grep" /tmp/result.txt; then
    echo "FAIL: Lines with 'grep' should be filtered out"
    exit 1
fi
if grep -q "\[" /tmp/result.txt; then
    echo "FAIL: Lines with '[' should be filtered out"
    exit 1
fi
if grep -q "nginx" /tmp/result.txt && grep -q "mysqld" /tmp/result.txt; then
    echo "PASS"
    exit 0
fi
echo "FAIL: Expected nginx and mysqld in output"
exit 1
