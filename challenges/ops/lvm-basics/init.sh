#!/bin/bash
apt-get update -qq && apt-get install -y -qq lvm2 > /dev/null 2>&1
rm -f /tmp/disk1.img /tmp/disk2.img /tmp/lvm_help.txt /tmp/lvm_concepts.txt
