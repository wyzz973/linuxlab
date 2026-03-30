#!/bin/bash
awk -F',' '{print $1"\t"$3}' /home/lab/employees.csv > /tmp/result.txt
