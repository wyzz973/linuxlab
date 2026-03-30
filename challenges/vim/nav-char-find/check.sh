#!/bin/bash
expected='name = "Alice"
city = "Beijing"
role = "admin"'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
