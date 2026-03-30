#!/bin/bash
dmesg | tail -20 > /tmp/result.txt 2>/dev/null || echo "dmesg not available" > /tmp/result.txt
