#!/bin/bash
if [ ! -f /tmp/passwd_sha256.txt ] || [ ! -s /tmp/passwd_sha256.txt ]; then
    echo "FAIL: /tmp/passwd_sha256.txt not found or empty"
    exit 1
fi
if [ ! -f /tmp/passwd_md5.txt ] || [ ! -s /tmp/passwd_md5.txt ]; then
    echo "FAIL: /tmp/passwd_md5.txt not found or empty"
    exit 1
fi
if [ ! -f /tmp/conf_checksums.txt ]; then
    echo "FAIL: /tmp/conf_checksums.txt not found"
    exit 1
fi
# Verify SHA256 format (64 hex chars)
if grep -qE "[a-f0-9]{64}" /tmp/passwd_sha256.txt; then
    echo "PASS"
    exit 0
fi
echo "FAIL: invalid SHA256 format"
exit 1
