#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "subdir1" /tmp/result.txt && grep -q "subdir2" /tmp/result.txt; then
    # Make sure regular files are not listed
    if grep -q "file1.txt\|file2.txt\|data.csv" /tmp/result.txt; then
        echo "FAIL: Regular files should not be in output"
        exit 1
    fi
    echo "PASS"
    exit 0
fi
echo "FAIL: Expected directory names in output"
exit 1
