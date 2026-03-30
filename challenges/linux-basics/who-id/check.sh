#!/bin/bash
if [ ! -f /tmp/id_result.txt ]; then
    echo "FAIL: /tmp/id_result.txt not found"
    exit 1
fi
if [ ! -f /tmp/whoami_result.txt ]; then
    echo "FAIL: /tmp/whoami_result.txt not found"
    exit 1
fi
if grep -q "uid=" /tmp/id_result.txt && grep -q "gid=" /tmp/id_result.txt; then
    echo "PASS"
    exit 0
else
    echo "FAIL: id output not correct"
    exit 1
fi
