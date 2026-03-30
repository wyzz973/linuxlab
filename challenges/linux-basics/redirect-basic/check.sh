#!/bin/bash
if [ ! -f /tmp/output.txt ]; then
    echo "FAIL: /tmp/output.txt not found"
    exit 1
fi
line1=$(sed -n '1p' /tmp/output.txt)
line2=$(sed -n '2p' /tmp/output.txt)
if [ "$line1" = "Hello Linux" ] && [ "$line2" = "Welcome to Shell" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected 'Hello Linux' then 'Welcome to Shell'"
    exit 1
fi
