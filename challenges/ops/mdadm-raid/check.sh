#!/bin/bash
if [ ! -f /tmp/mdadm_help.txt ]; then
    echo "FAIL: /tmp/mdadm_help.txt not found"
    exit 1
fi
if [ ! -f /tmp/raid_status.txt ]; then
    echo "FAIL: /tmp/raid_status.txt not found"
    exit 1
fi
if [ ! -f /tmp/raid_levels.txt ] || [ ! -s /tmp/raid_levels.txt ]; then
    echo "FAIL: /tmp/raid_levels.txt not found or empty"
    exit 1
fi
if grep -qiE "RAID 0|RAID 1|RAID 5|RAID.10|raid0|raid1" /tmp/raid_levels.txt; then
    echo "PASS"
    exit 0
fi
echo "FAIL: RAID levels not properly documented"
exit 1
