#!/bin/bash
if [ ! -f /tmp/index.html ]; then
    echo "FAIL: /tmp/index.html not found"
    exit 1
fi
if [ ! -s /tmp/index.html ]; then
    echo "FAIL: /tmp/index.html is empty"
    exit 1
fi
if [ ! -f /tmp/wget.log ]; then
    echo "FAIL: /tmp/wget.log not found"
    exit 1
fi
if [ ! -f /tmp/wget_recursive_help.txt ]; then
    echo "FAIL: /tmp/wget_recursive_help.txt not found"
    exit 1
fi
echo "PASS"
exit 0
