#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
lines=$(wc -l < /tmp/result.txt | tr -d ' ')
if [ "$lines" -eq 3 ]; then
    # Check that delta (largest) appears first
    first_line=$(head -1 /tmp/result.txt)
    if echo "$first_line" | grep -q "delta"; then
        echo "PASS"
        exit 0
    fi
fi
echo "FAIL: Expected 3 directories sorted by size (delta should be first)"
exit 1
