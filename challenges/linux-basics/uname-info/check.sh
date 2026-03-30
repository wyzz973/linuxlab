#!/bin/bash
if [ ! -f /tmp/uname_result.txt ] || [ ! -f /tmp/hostname_result.txt ]; then
    echo "FAIL: Result files not found"
    exit 1
fi
if grep -qi "linux\|darwin\|unix" /tmp/uname_result.txt && [ -s /tmp/hostname_result.txt ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected system info in output"
    exit 1
fi
