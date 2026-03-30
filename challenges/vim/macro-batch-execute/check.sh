#!/bin/bash
expected='"name": "Alice",
"age": "30",
"city": "Shanghai",
"role": "engineer",
"level": "senior",
"team": "backend",'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
