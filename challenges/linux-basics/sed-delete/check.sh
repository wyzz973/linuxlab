#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
# No comment lines
if grep -q "^#" /tmp/result.txt; then
    echo "FAIL: Comment lines still present"
    exit 1
fi
# No empty lines
if grep -q "^$" /tmp/result.txt; then
    echo "FAIL: Empty lines still present"
    exit 1
fi
# Should still have code
if grep -q "def calculate" /tmp/result.txt; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Code lines missing"
    exit 1
fi
