#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
# Check long format is used (should contain permission strings)
if ! grep -qE "^[-drwx]" /tmp/result.txt; then
    echo "FAIL: Output doesn't appear to be in long format"
    exit 1
fi
# Check that large_file appears before small_file (sorted by size)
large_line=$(grep -n "large_file" /tmp/result.txt | head -1 | cut -d: -f1)
small_line=$(grep -n "small_file" /tmp/result.txt | head -1 | cut -d: -f1)
if [ -n "$large_line" ] && [ -n "$small_line" ] && [ "$large_line" -lt "$small_line" ]; then
    echo "PASS: Files sorted by size correctly"
    exit 0
else
    echo "FAIL: Files not sorted by size (large to small)"
    exit 1
fi
