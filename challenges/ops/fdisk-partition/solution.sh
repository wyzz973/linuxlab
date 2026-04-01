#!/bin/bash
dd if=/dev/zero of=/tmp/partition_disk.img bs=1M count=200 2>/dev/null
fdisk -l /tmp/partition_disk.img > /tmp/fdisk_info.txt 2>&1
lsblk > /tmp/block_devices.txt 2>&1 || ls -la /dev/sd* /dev/vd* > /tmp/block_devices.txt 2>&1 || echo "No block devices" > /tmp/block_devices.txt
