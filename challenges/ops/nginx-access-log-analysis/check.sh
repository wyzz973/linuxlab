#!/bin/bash
if [ ! -f /tmp/top_ips.txt ] || [ ! -s /tmp/top_ips.txt ]; then
    echo "FAIL: /tmp/top_ips.txt not found or empty"
    exit 1
fi
if [ ! -f /tmp/status_codes.txt ] || [ ! -s /tmp/status_codes.txt ]; then
    echo "FAIL: /tmp/status_codes.txt not found or empty"
    exit 1
fi
if [ ! -f /tmp/not_found.txt ] || [ ! -s /tmp/not_found.txt ]; then
    echo "FAIL: /tmp/not_found.txt not found or empty"
    exit 1
fi
# Verify top_ips has IP addresses
if ! grep -qE "[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+" /tmp/top_ips.txt; then
    echo "FAIL: /tmp/top_ips.txt does not contain IP addresses"
    exit 1
fi
# Verify status_codes has HTTP status codes
if ! grep -qE "[0-9]+ +[0-9]{3}" /tmp/status_codes.txt; then
    echo "FAIL: /tmp/status_codes.txt does not contain status codes"
    exit 1
fi
# Verify not_found has URL paths
if ! grep -qE "^/" /tmp/not_found.txt; then
    echo "FAIL: /tmp/not_found.txt does not contain URL paths"
    exit 1
fi
echo "PASS"
exit 0
