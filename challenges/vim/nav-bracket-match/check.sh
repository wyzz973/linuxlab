#!/bin/bash
expected='function main() {
    let data = [1, 2, 3];
    let result = process(data);
    console.log(result);
}'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
