#!/bin/bash
mkdir -p /home/lab/{alpha,beta,gamma,delta,epsilon}
dd if=/dev/zero of=/home/lab/alpha/data.bin bs=1024 count=500 2>/dev/null
dd if=/dev/zero of=/home/lab/beta/data.bin bs=1024 count=1000 2>/dev/null
dd if=/dev/zero of=/home/lab/gamma/data.bin bs=1024 count=200 2>/dev/null
dd if=/dev/zero of=/home/lab/delta/data.bin bs=1024 count=2000 2>/dev/null
dd if=/dev/zero of=/home/lab/epsilon/data.bin bs=1024 count=100 2>/dev/null
