#!/bin/bash
expected='// int x = 10;
// int y = 20;
// int z = 30;
// int w = 40;
// int v = 50;'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
