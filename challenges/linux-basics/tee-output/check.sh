#!/bin/bash
if [ ! -f /tmp/passed.txt ]; then
    echo "FAIL: /tmp/passed.txt not found"
    exit 1
fi
pass_count=$(grep -c "PASS" /tmp/passed.txt)
fail_count=$(grep -c "FAIL" /tmp/passed.txt 2>/dev/null) || fail_count=0
if [ "$pass_count" -eq 4 ] && [ "$fail_count" -eq 0 ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected 4 PASS lines and 0 FAIL lines"
    exit 1
fi
