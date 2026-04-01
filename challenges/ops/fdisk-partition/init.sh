#!/bin/bash
apt-get update -qq && apt-get install -y -qq fdisk util-linux > /dev/null 2>&1
rm -f /tmp/partition_disk.img /tmp/fdisk_info.txt /tmp/block_devices.txt
