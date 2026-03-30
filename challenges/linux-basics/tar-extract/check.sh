#!/bin/bash
if [ ! -f /home/lab/extracted/data/users.csv ]; then
    echo "FAIL: /home/lab/extracted/data/users.csv not found"
    exit 1
fi
if grep -q "Alice" /home/lab/extracted/data/users.csv; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Extracted file content incorrect"
    exit 1
fi
