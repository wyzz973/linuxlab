#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
# Check that names and salaries are present, but department is not
if grep -q "Alice" /tmp/result.txt && grep -q "95000" /tmp/result.txt; then
    if grep -q "Engineering" /tmp/result.txt || grep -q "Marketing" /tmp/result.txt; then
        echo "FAIL: Department column should not be in output"
        exit 1
    fi
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected name and salary columns"
    exit 1
fi
