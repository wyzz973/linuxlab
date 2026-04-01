#!/bin/bash
rm -f /tmp/top_cpu.txt /tmp/top_mem.txt /tmp/proc_count.txt
# Create some CPU load
(dd if=/dev/urandom bs=1M count=10 of=/dev/null 2>/dev/null) &
