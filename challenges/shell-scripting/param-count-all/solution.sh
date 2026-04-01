#!/bin/bash
echo "Count: $#"
for arg in "$@"; do
    echo "$arg"
done
result=""
for arg in "$@"; do
    if [ -n "$result" ]; then
        result="$result,$arg"
    else
        result="$arg"
    fi
done
echo "All: $result"
