#!/bin/bash
declare -A capitals
capitals[China]="Beijing"
capitals[Japan]="Tokyo"
capitals[France]="Paris"

echo "${capitals[China]}"

for key in $(echo "${!capitals[@]}" | tr ' ' '\n' | sort); do
    echo "$key"
done

for key in $(echo "${!capitals[@]}" | tr ' ' '\n' | sort); do
    echo "${capitals[$key]}"
done

echo "${#capitals[@]}"
