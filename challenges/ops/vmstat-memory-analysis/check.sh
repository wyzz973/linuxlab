#!/bin/bash
if [ ! -f /tmp/vmstat_output.txt ] || [ ! -s /tmp/vmstat_output.txt ]; then
    echo "FAIL: /tmp/vmstat_output.txt not found or empty"
    exit 1
fi
if [ ! -f /tmp/mem_info.txt ] || [ ! -s /tmp/mem_info.txt ]; then
    echo "FAIL: /tmp/mem_info.txt not found or empty"
    exit 1
fi
if [ ! -f /tmp/swap_info.txt ]; then
    echo "FAIL: /tmp/swap_info.txt not found"
    exit 1
fi
if grep -qE "procs|memory|swpd|free|buff" /tmp/vmstat_output.txt; then
    echo "PASS"
    exit 0
fi
echo "FAIL: vmstat output format incorrect"
exit 1
