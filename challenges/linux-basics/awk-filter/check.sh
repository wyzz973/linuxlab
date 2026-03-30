#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "Laptop" /tmp/result.txt && grep -q "6000" /tmp/result.txt; then
    if grep -q "Mouse" /tmp/result.txt && grep -q "USB_Cable" /tmp/result.txt; then
        # Mouse=2500, USB=2000 also > 1000
        echo "PASS"
        exit 0
    fi
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected Laptop with total 6000"
    exit 1
fi
