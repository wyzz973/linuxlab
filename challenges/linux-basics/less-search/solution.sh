#!/bin/bash
grep -C 2 'CRITICAL' /home/lab/large_log.txt > /tmp/result.txt
