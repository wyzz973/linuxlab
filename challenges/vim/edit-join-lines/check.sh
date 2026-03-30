#!/bin/bash
expected='SELECT name, email, age FROM users WHERE active = true ORDER BY name;'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
