#!/bin/bash
iostat -x 1 3 > /tmp/iostat_output.txt 2>&1 || iostat 1 3 > /tmp/iostat_output.txt 2>&1
cat /proc/diskstats > /tmp/diskstats.txt 2>&1
df -h > /tmp/disk_usage.txt
