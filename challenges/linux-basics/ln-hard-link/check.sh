#!/bin/bash
if [ ! -f /home/lab/hardlink.txt ]; then
    echo "FAIL: /home/lab/hardlink.txt not found"
    exit 1
fi
# Check it's a hard link (same inode)
inode1=$(stat -c %i /home/lab/original.txt 2>/dev/null || stat -f %i /home/lab/original.txt 2>/dev/null)
inode2=$(stat -c %i /home/lab/hardlink.txt 2>/dev/null || stat -f %i /home/lab/hardlink.txt 2>/dev/null)
if [ "$inode1" = "$inode2" ]; then
    echo "PASS: Hard link created (same inode: $inode1)"
    exit 0
else
    echo "FAIL: Not a hard link (different inodes: $inode1 vs $inode2)"
    exit 1
fi
