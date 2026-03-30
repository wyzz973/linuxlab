#!/bin/bash
mkdir -p /home/lab/storage
dd if=/dev/zero of=/home/lab/storage/bigfile1.dat bs=1024 count=100 2>/dev/null
dd if=/dev/zero of=/home/lab/storage/bigfile2.dat bs=1024 count=200 2>/dev/null
dd if=/dev/zero of=/home/lab/storage/small1.txt bs=1024 count=5 2>/dev/null
dd if=/dev/zero of=/home/lab/storage/small2.txt bs=1024 count=10 2>/dev/null
echo "tiny" > /home/lab/storage/tiny.txt
