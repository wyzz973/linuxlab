#!/bin/bash
sort /home/lab/access.log | uniq -c | sort -rn > /tmp/result.txt
