#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
first_line=$(head -1 /tmp/result.txt)
if echo "$first_line" | grep -q "v3.0"; then
    echo "PASS"
    exit 0
else
    echo "FAIL: First line should contain v3.0 (latest entry)"
    exit 1
fi
