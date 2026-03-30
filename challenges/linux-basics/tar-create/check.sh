#!/bin/bash
if [ ! -f /home/lab/website-backup.tar.gz ]; then
    echo "FAIL: /home/lab/website-backup.tar.gz not found"
    exit 1
fi
# Verify it's a valid tar.gz and contains expected files
if tar tzf /home/lab/website-backup.tar.gz 2>/dev/null | grep -q "index.html"; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Archive doesn't contain expected files"
    exit 1
fi
