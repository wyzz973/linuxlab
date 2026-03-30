#!/bin/bash
if [ ! -f /tmp/installed.txt ]; then
    echo "FAIL: /tmp/installed.txt not found"
    exit 1
fi
if [ -s /tmp/installed.txt ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Installed packages list is empty"
    exit 1
fi
