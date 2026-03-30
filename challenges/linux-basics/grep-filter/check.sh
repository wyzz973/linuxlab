#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
error_count=$(grep -c "ERROR" /tmp/result.txt)
total_lines=$(wc -l < /tmp/result.txt | tr -d ' ')
if [ "$error_count" -eq 3 ] && [ "$total_lines" -eq 3 ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected 3 ERROR lines, got $error_count errors in $total_lines lines"
    exit 1
fi
