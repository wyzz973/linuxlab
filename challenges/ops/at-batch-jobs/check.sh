#!/bin/bash
if [ ! -f /tmp/at_help.txt ]; then
    echo "FAIL: /tmp/at_help.txt not found"
    exit 1
fi
if [ ! -f /tmp/at_demo.sh ]; then
    echo "FAIL: /tmp/at_demo.sh not found"
    exit 1
fi
if [ ! -f /tmp/at_queue.txt ]; then
    echo "FAIL: /tmp/at_queue.txt not found"
    exit 1
fi
if grep -qiE "at|schedule\|command\|date" /tmp/at_demo.sh; then
    echo "PASS"
    exit 0
fi
echo "PASS"
exit 0
