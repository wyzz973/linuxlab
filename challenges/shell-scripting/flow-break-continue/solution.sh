#!/bin/bash
line=""
for i in $(seq 1 20); do
    if [ $((i % 3)) -eq 0 ]; then
        continue
    fi
    if [ -n "$line" ]; then
        line="$line $i"
    else
        line="$i"
    fi
done
echo "$line"

sum=0
for i in $(seq 1 100); do
    sum=$((sum + i))
    if [ $sum -gt 200 ]; then
        echo "Stopped at $i, sum=$sum"
        break
    fi
done
