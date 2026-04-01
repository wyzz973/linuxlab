#!/bin/bash
ps aux --sort=-%cpu | head -11 > /tmp/top_cpu.txt
ps aux --sort=-%mem | head -11 > /tmp/top_mem.txt
ps aux | tail -n +2 | wc -l | tr -d " " > /tmp/proc_count.txt
