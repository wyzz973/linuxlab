#!/bin/bash
sar 1 3 > /tmp/sar_cpu.txt 2>&1
sar -r 1 3 > /tmp/sar_mem.txt 2>&1
sar -n DEV 1 3 > /tmp/sar_net.txt 2>&1
