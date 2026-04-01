#!/bin/bash
apt-get update -qq && apt-get install -y -qq procps > /dev/null 2>&1
rm -f /tmp/vmstat_output.txt /tmp/mem_info.txt /tmp/swap_info.txt
