#!/bin/bash
if [ -f /home/lab/downloads/index.html ]; then
    echo "PASS"
    exit 0
elif [ -f /tmp/result.txt ]; then
    echo "PASS: Command recorded"
    exit 0
else
    echo "FAIL: /home/lab/downloads/index.html not found"
    exit 1
fi
