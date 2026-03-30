#!/bin/bash
if [ ! -f /tmp/env_result.txt ]; then
    echo "FAIL: /tmp/env_result.txt not found"
    exit 1
fi
if [ ! -f /tmp/path_result.txt ]; then
    echo "FAIL: /tmp/path_result.txt not found"
    exit 1
fi
if grep -q "MY_APP_ENV=production" /tmp/env_result.txt; then
    if grep -q "/" /tmp/path_result.txt; then
        echo "PASS"
        exit 0
    fi
fi
echo "FAIL: MY_APP_ENV not found or PATH not saved"
exit 1
