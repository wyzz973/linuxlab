#!/bin/bash
awk '{total=$2*$3; if(total>1000) print $1, total}' /home/lab/sales.txt > /tmp/result.txt
