#!/bin/bash
mkdir -p /home/lab/data
dd if=/dev/zero of=/home/lab/data/large_file.dat bs=1024 count=100 2>/dev/null
dd if=/dev/zero of=/home/lab/data/medium_file.dat bs=1024 count=50 2>/dev/null
dd if=/dev/zero of=/home/lab/data/small_file.dat bs=1024 count=10 2>/dev/null
echo "tiny" > /home/lab/data/tiny_file.txt
