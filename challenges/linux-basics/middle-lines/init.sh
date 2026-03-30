#!/bin/bash
mkdir -p /home/lab
for i in $(seq 1 50); do
    echo "Line $i: Data entry number $i with value $((RANDOM % 1000))"
done > /home/lab/data.txt
