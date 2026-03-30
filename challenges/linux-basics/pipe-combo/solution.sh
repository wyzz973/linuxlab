#!/bin/bash
awk '{print $1}' /home/lab/web_access.log | sort | uniq -c | sort -rn | head -5 > /tmp/result.txt
