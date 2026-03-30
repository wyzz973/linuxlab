#!/bin/bash
if which tree &>/dev/null; then
    if [ -f /tmp/result.txt ]; then
        echo "PASS"
        exit 0
    fi
    echo "PASS: tree is installed"
    exit 0
else
    echo "FAIL: tree is not installed"
    exit 1
fi
