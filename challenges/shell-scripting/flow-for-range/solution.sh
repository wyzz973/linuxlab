#!/bin/bash
echo {1..5}

for i in $(seq 2 2 10); do
    echo "$i"
done

sum=0
for i in $(seq 1 100); do
    sum=$((sum + i))
done
echo "$sum"
