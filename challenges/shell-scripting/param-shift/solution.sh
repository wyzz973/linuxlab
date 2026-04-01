#!/bin/bash
count=1
while [ $# -gt 0 ]; do
    remaining=$(($# - 1))
    echo "Processing #$count: $1 ($remaining remaining)"
    shift
    count=$((count + 1))
done
