#!/bin/bash
mkdir -p /home/lab/{logs,data,config,cache}
dd if=/dev/zero of=/home/lab/logs/app.log bs=1024 count=500 2>/dev/null
dd if=/dev/zero of=/home/lab/data/database.db bs=1024 count=1000 2>/dev/null
echo "small config" > /home/lab/config/app.conf
dd if=/dev/zero of=/home/lab/cache/temp.dat bs=1024 count=200 2>/dev/null
