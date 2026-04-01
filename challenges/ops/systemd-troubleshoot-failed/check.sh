#!/bin/bash
if [ ! -f /tmp/broken_status.txt ]; then
    echo "FAIL: /tmp/broken_status.txt not found"
    exit 1
fi
if [ ! -f /tmp/broken_logs.txt ]; then
    echo "FAIL: /tmp/broken_logs.txt not found"
    exit 1
fi
if [ ! -f /tmp/fixed_service.txt ]; then
    echo "FAIL: /tmp/fixed_service.txt not found"
    exit 1
fi
# Check that the fix references a valid script
if grep -q "fixed_app.sh" /tmp/fixed_service.txt; then
    echo "PASS"
    exit 0
fi
# Also accept if ExecStart points to an existing file
EXEC_PATH=$(grep ExecStart /tmp/fixed_service.txt | cut -d= -f2 | tr -d " ")
if [ -f "$EXEC_PATH" ] && [ -x "$EXEC_PATH" ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: service not properly fixed"
exit 1
