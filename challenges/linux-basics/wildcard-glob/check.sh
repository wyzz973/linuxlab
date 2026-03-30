#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "data1.txt" /tmp/result.txt && grep -q "data_final.txt" /tmp/result.txt; then
    if grep -q "report.txt" /tmp/result.txt; then
        echo "FAIL: report.txt should not be included"
        exit 1
    fi
    if grep -q "data1.csv" /tmp/result.txt; then
        echo "FAIL: data1.csv should not be included"
        exit 1
    fi
    echo "PASS"
    exit 0
fi
echo "FAIL: Expected data*.txt files"
exit 1
