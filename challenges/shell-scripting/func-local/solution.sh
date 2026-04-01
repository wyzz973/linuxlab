#!/bin/bash
count=0

increment() {
    local local_count=100
    count=$((count + 1))
    echo "Inside: count=$count, local_count=$local_count"
}

increment
increment
increment
echo "Outside: count=$count"
echo "Outside: local_count=$local_count"
