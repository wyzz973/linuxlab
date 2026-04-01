#!/bin/bash
colors=(red green blue yellow purple)
echo "${#colors[@]}"
for color in "${colors[@]}"; do
    echo "$color"
done
for i in $(seq 0 $((${#colors[@]}-1))); do
    echo "$i: ${colors[$i]}"
done
