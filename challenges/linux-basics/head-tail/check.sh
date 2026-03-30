#!/bin/bash
if [ ! -f /tmp/head_result.txt ] || [ ! -f /tmp/tail_result.txt ]; then
    echo "FAIL: Result files not found"
    exit 1
fi
head_lines=$(wc -l < /tmp/head_result.txt | tr -d ' ')
tail_lines=$(wc -l < /tmp/tail_result.txt | tr -d ' ')
if [ "$head_lines" -eq 5 ] && [ "$tail_lines" -eq 3 ]; then
    if grep -q "Server started" /tmp/head_result.txt && grep -q "Connection restored" /tmp/tail_result.txt; then
        echo "PASS"
        exit 0
    fi
fi
echo "FAIL: Incorrect number of lines or content"
exit 1
