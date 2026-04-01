#!/bin/bash
scores=(85 92 78 95 88 76 90 82 97 70)
count=${#scores[@]}
max=${scores[0]}
min=${scores[0]}
total=0

for score in "${scores[@]}"; do
    total=$((total + score))
    if [ "$score" -gt "$max" ]; then
        max=$score
    fi
    if [ "$score" -lt "$min" ]; then
        min=$score
    fi
done

average=$((total / count))

echo "Students: $count"
echo "Highest: $max"
echo "Lowest: $min"
echo "Total: $total"
echo "Average: $average"
