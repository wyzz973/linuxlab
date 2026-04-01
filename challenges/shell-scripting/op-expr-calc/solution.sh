#!/bin/bash
while read a op b; do
    result=$(( a $op b ))
    echo "$a $op $b = $result"
done < /home/learner/calculations.txt
