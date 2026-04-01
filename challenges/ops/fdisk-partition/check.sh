#!/bin/bash
if [ ! -f /tmp/partition_disk.img ]; then
    echo "FAIL: /tmp/partition_disk.img not found"
    exit 1
fi
SIZE=$(stat -c%s /tmp/partition_disk.img 2>/dev/null || stat -f%z /tmp/partition_disk.img 2>/dev/null)
if [ "$SIZE" -lt 190000000 ]; then
    echo "FAIL: disk image too small"
    exit 1
fi
if [ ! -f /tmp/fdisk_info.txt ]; then
    echo "FAIL: /tmp/fdisk_info.txt not found"
    exit 1
fi
if [ ! -f /tmp/block_devices.txt ]; then
    echo "FAIL: /tmp/block_devices.txt not found"
    exit 1
fi
echo "PASS"
exit 0
