#!/bin/bash
apt-get update -qq && apt-get install -y -qq sysstat > /dev/null 2>&1
rm -f /tmp/sar_cpu.txt /tmp/sar_mem.txt /tmp/sar_net.txt
