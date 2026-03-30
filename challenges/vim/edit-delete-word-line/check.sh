#!/bin/bash
expected='int count = 0;
count += 1;
printf("%d", count);'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
