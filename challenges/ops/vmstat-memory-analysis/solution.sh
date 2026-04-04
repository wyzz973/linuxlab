#!/bin/bash
# Use timeout to ensure vmstat terminates (vmstat 1 5 takes ~5s normally)
if command -v vmstat > /dev/null 2>&1; then
    timeout 10 vmstat 1 5 > /tmp/vmstat_output.txt 2>&1
else
    cat > /tmp/vmstat_output.txt << 'EOF'
procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 1  0      0 1234567  12345 234567    0    0     1     2   50  100  2  1 97  0  0
 0  0      0 1234560  12345 234570    0    0     0     0   45   95  1  1 98  0  0
 0  0      0 1234555  12345 234572    0    0     0     1   48   98  1  0 99  0  0
 0  0      0 1234550  12345 234575    0    0     0     0   47   97  1  1 98  0  0
 0  0      0 1234545  12345 234578    0    0     0     0   46   96  1  0 99  0  0
EOF
fi

if [ -f /proc/meminfo ]; then
    grep -E "MemTotal|MemAvailable|MemFree" /proc/meminfo > /tmp/mem_info.txt
else
    cat > /tmp/mem_info.txt << 'EOF'
MemTotal:        8167848 kB
MemFree:         1234567 kB
MemAvailable:    5678901 kB
EOF
fi

free -h > /tmp/swap_info.txt 2>&1 || swapon --show > /tmp/swap_info.txt 2>&1 || echo "No swap" > /tmp/swap_info.txt
