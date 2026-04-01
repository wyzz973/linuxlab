#!/bin/bash
vmstat 1 5 > /tmp/vmstat_output.txt 2>&1
grep -E "MemTotal|MemAvailable|MemFree" /proc/meminfo > /tmp/mem_info.txt
free -h > /tmp/swap_info.txt 2>&1 || swapon --show > /tmp/swap_info.txt 2>&1 || echo "No swap" > /tmp/swap_info.txt
