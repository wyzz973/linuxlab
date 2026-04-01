#!/bin/bash
factorial() {
    if [ "$1" -le 1 ]; then
        echo 1
    else
        local prev=$(factorial $(($1 - 1)))
        echo $(($1 * prev))
    fi
}

fibonacci() {
    if [ "$1" -eq 0 ]; then
        echo 0
    elif [ "$1" -eq 1 ]; then
        echo 1
    else
        local a=$(fibonacci $(($1 - 1)))
        local b=$(fibonacci $(($1 - 2)))
        echo $((a + b))
    fi
}

result=$(factorial 6)
echo "6! = $result"

line=""
for i in $(seq 0 9); do
    val=$(fibonacci $i)
    if [ -n "$line" ]; then
        line="$line $val"
    else
        line="$val"
    fi
done
echo "$line"
