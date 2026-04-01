#!/bin/bash
for os in Linux macOS Windows Android iOS; do
    echo "$os"
done

count=1
while read item; do
    echo "$count. $item"
    count=$((count + 1))
done < /home/learner/items.txt
