#!/bin/bash
if [ ! -f /home/lab/combined.txt ]; then
    echo "FAIL: /home/lab/combined.txt not found"
    exit 1
fi
if grep -q "Part 1 content" /home/lab/combined.txt && grep -q "Part 2 content" /home/lab/combined.txt; then
    # Check order
    line1=$(grep -n "Part 1" /home/lab/combined.txt | head -1 | cut -d: -f1)
    line2=$(grep -n "Part 2" /home/lab/combined.txt | head -1 | cut -d: -f1)
    if [ "$line1" -lt "$line2" ]; then
        echo "PASS"
        exit 0
    fi
fi
echo "FAIL: Content not merged correctly"
exit 1
