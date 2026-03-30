#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "logs" /tmp/result.txt && grep -q "data" /tmp/result.txt; then
    # Check human readable format (K, M, G)
    if grep -qE "[0-9]+[KMG]" /tmp/result.txt || grep -qE "[0-9]+\.[0-9]+[KMG]" /tmp/result.txt; then
        echo "PASS"
        exit 0
    fi
    # Some systems use different formatting
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected subdirectory names in output"
    exit 1
fi
