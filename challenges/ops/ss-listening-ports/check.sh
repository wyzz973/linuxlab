#!/bin/bash
if [ ! -f /tmp/listening.txt ]; then
    echo "FAIL: /tmp/listening.txt not found"
    exit 1
fi
if [ ! -f /tmp/ports.txt ]; then
    echo "FAIL: /tmp/ports.txt not found"
    exit 1
fi
if ! grep -q "LISTEN\|State" /tmp/listening.txt; then
    echo "FAIL: /tmp/listening.txt does not contain listening sockets"
    exit 1
fi
if [ ! -s /tmp/ports.txt ]; then
    echo "FAIL: /tmp/ports.txt is empty"
    exit 1
fi
# Check that at least one known port is present
if grep -qE "8080|9090|3333" /tmp/ports.txt; then
    echo "PASS"
    exit 0
fi
echo "FAIL: expected ports not found in /tmp/ports.txt"
exit 1
