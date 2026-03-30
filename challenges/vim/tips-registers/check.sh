#!/bin/bash
expected='# Source data (do not modify above lines)
alpha [2]
beta [3]
gamma [1]

# OUTPUT (paste in order below):
gamma [1]
alpha [2]
beta [3]'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
