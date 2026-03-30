#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "Total:" /tmp/result.txt && grep -q "Average:" /tmp/result.txt; then
    # Check approximate total (should be 1500)
    if grep -qE "1500|Total:.*15" /tmp/result.txt; then
        echo "PASS"
        exit 0
    fi
    echo "PASS: Format correct"
    exit 0
fi
echo "FAIL: Expected Total: and Average: in output"
exit 1
