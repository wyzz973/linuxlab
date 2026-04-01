#!/bin/bash
n=5
result=1
while [ $n -gt 0 ]; do
    result=$((result * n))
    n=$((n - 1))
done
echo "5! = $result"

count=5
while [ $count -gt 0 ]; do
    echo "$count"
    count=$((count - 1))
done
echo "Go!"
