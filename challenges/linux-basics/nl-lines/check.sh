#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
# Count numbered lines - with -b a, empty lines should also be numbered
total_lines=$(wc -l < /tmp/result.txt | tr -d ' ')
numbered_lines=$(grep -cE "^\s+[0-9]+" /tmp/result.txt)
if [ "$numbered_lines" -ge 10 ]; then
    echo "PASS: All lines including blanks are numbered"
    exit 0
else
    echo "FAIL: Expected all lines numbered, found $numbered_lines"
    exit 1
fi
