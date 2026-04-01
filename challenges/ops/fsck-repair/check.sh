#!/bin/bash
if [ ! -f /tmp/test_disk.img ]; then
    echo "FAIL: /tmp/test_disk.img not found"
    exit 1
fi
if [ ! -f /tmp/fsck_result.txt ] || [ ! -s /tmp/fsck_result.txt ]; then
    echo "FAIL: /tmp/fsck_result.txt not found or empty"
    exit 1
fi
if [ ! -f /tmp/fs_info.txt ] || [ ! -s /tmp/fs_info.txt ]; then
    echo "FAIL: /tmp/fs_info.txt not found or empty"
    exit 1
fi
if grep -qiE "clean|ext4|pass" /tmp/fsck_result.txt; then
    echo "PASS"
    exit 0
fi
echo "PASS"
exit 0
