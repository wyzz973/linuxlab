#!/bin/bash
mkdir -p /home/lab/exists
echo "file1" > /home/lab/exists/file1.txt
rm -rf /home/lab/nonexistent
rm -f /tmp/stdout.txt /tmp/stderr.txt
