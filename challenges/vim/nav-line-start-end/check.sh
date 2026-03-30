#!/bin/bash
expected='int x = 10;
int y = 20;
int z = x + y;
printf("%d", z);
return 0;'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
