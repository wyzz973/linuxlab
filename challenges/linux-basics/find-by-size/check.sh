#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "bigfile1" /tmp/result.txt && grep -q "bigfile2" /tmp/result.txt; then
    if grep -q "small1" /tmp/result.txt || grep -q "tiny" /tmp/result.txt; then
        echo "FAIL: Small files should not be included"
        exit 1
    fi
    echo "PASS"
    exit 0
else
    echo "FAIL: Large files not found in result"
    exit 1
fi
