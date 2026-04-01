#!/bin/bash
if [ ! -f /tmp/lvm_help.txt ]; then
    echo "FAIL: /tmp/lvm_help.txt not found"
    exit 1
fi
if [ ! -f /tmp/disk1.img ]; then
    echo "FAIL: /tmp/disk1.img not found"
    exit 1
fi
if [ ! -f /tmp/disk2.img ]; then
    echo "FAIL: /tmp/disk2.img not found"
    exit 1
fi
SIZE1=$(stat -c%s /tmp/disk1.img 2>/dev/null || stat -f%z /tmp/disk1.img 2>/dev/null)
if [ "$SIZE1" -lt 90000000 ]; then
    echo "FAIL: disk1.img is too small (should be ~100MB)"
    exit 1
fi
if [ ! -f /tmp/lvm_concepts.txt ] || [ ! -s /tmp/lvm_concepts.txt ]; then
    echo "FAIL: /tmp/lvm_concepts.txt not found or empty"
    exit 1
fi
echo "PASS"
exit 0
