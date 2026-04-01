#!/bin/bash
for((i=1; i<=5; i++)); do
    echo "$i^2 = $((i * i))"
done

line=""
for((j=1; j<=5; j++)); do
    if [ -n "$line" ]; then
        line="$line "
    fi
    line="${line}5x${j}=$((5 * j))"
done
echo "$line"
