#!/bin/bash
dd if=/dev/zero of=/tmp/test_disk.img bs=1M count=50 2>/dev/null
mkfs.ext4 -F /tmp/test_disk.img 2>/dev/null
fsck.ext4 -n /tmp/test_disk.img > /tmp/fsck_result.txt 2>&1
dumpe2fs /tmp/test_disk.img > /tmp/fs_info.txt 2>&1 || tune2fs -l /tmp/test_disk.img > /tmp/fs_info.txt 2>&1
