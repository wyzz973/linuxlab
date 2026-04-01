#!/bin/bash
sum=0
n=0
until [ $sum -gt 50 ]; do
    n=$((n + 1))
    sum=$((sum + n))
done
echo "Sum: $sum, stopped at: $n"

val=100
until [ $val -le 0 ]; do
    echo "$val"
    val=$((val - 7))
done
