#!/bin/bash
awk '{sum+=$1} END{print "Total:", sum, "Average:", sum/NR}' /home/lab/transactions.txt > /tmp/result.txt
