#!/bin/bash
is_even() {
    if [ $(($1 % 2)) -eq 0 ]; then
        return 0
    else
        return 1
    fi
}

max() {
    if [ "$1" -gt "$2" ]; then
        echo "$1"
    else
        echo "$2"
    fi
}

is_even 4 && echo "4 is even" || echo "4 is odd"
is_even 7 && echo "7 is even" || echo "7 is odd"
result=$(max 15 23)
echo "Max of 15 and 23 is: $result"
