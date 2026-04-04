#!/bin/bash
# Ensure procps (vmstat, free) is available
if ! command -v vmstat > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq procps > /dev/null 2>&1 || true
fi
rm -f /tmp/vmstat_output.txt /tmp/mem_info.txt /tmp/swap_info.txt
