#!/bin/bash
if [ ! -f /tmp/tail_output.txt ] || [ ! -s /tmp/tail_output.txt ]; then
    echo "FAIL: /tmp/tail_output.txt not found or empty"
    exit 1
fi
TAIL_LINES=$(wc -l < /tmp/tail_output.txt | tr -d " ")
if [ "$TAIL_LINES" -ne 20 ]; then
    echo "FAIL: expected 20 lines in tail output, got $TAIL_LINES"
    exit 1
fi
if [ ! -f /tmp/live_errors.txt ]; then
    echo "FAIL: /tmp/live_errors.txt not found"
    exit 1
fi
if [ ! -f /tmp/line_count.txt ] || [ ! -s /tmp/line_count.txt ]; then
    echo "FAIL: /tmp/line_count.txt not found or empty"
    exit 1
fi
echo "PASS"
exit 0
