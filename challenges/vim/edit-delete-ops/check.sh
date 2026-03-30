#!/bin/bash
expected='name=production
debug=false
port=8080
log_level=warn'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
