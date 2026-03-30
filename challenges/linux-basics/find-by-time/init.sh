#!/bin/bash
mkdir -p /home/lab/logs
echo "recent error log" > /home/lab/logs/app.log
echo "recent access log" > /home/lab/logs/access.log
echo "recent debug" > /home/lab/logs/debug.log
# Create old files
echo "old log" > /home/lab/logs/old_app.log
touch -t 202401010000 /home/lab/logs/old_app.log
echo "old data" > /home/lab/logs/old_data.log
touch -t 202401010000 /home/lab/logs/old_data.log
echo "not a log" > /home/lab/logs/readme.txt
