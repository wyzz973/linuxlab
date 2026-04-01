#!/bin/bash
count_with_star() {
    local count=0
    for arg in "$*"; do
        count=$((count + 1))
    done
    echo "Using \$*: $count iterations"
}

count_with_at() {
    local count=0
    for arg in "$@"; do
        count=$((count + 1))
    done
    echo "Using \$@: $count iterations"
}

count_with_star "$@"
count_with_at "$@"
