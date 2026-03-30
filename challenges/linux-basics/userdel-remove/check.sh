#!/bin/bash
if id exemployee &>/dev/null; then
    echo "FAIL: User exemployee still exists"
    exit 1
fi
if [ -d /home/exemployee ]; then
    echo "FAIL: Home directory still exists"
    exit 1
fi
echo "PASS"
exit 0
