#!/bin/bash
apt-get update -qq && apt-get install -y -qq e2fsprogs > /dev/null 2>&1
rm -f /tmp/test_disk.img /tmp/fsck_result.txt /tmp/fs_info.txt
