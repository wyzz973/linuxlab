#!/bin/bash
find /home/lab/logs -name '*.log' -mtime -7 > /tmp/result.txt
