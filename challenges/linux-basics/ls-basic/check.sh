#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q ".hidden_config" /tmp/result.txt && grep -q ".bashrc" /tmp/result.txt; then
    echo "PASS: Hidden files are listed"
    exit 0
else
    echo "FAIL: Hidden files not found in output"
    exit 1
fi
