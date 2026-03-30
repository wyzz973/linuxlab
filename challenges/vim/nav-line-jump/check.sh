#!/bin/bash
expected='[server]
host = 0.0.0.0
port = 8080
workers = 4

[logging]
debug = true
level = info
file = /var/log/app.log
rotate = daily'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
