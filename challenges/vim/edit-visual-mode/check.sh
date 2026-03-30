#!/bin/bash
expected='function calculate(a, b) {
    let sum = a + b;
    return sum;
}'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
