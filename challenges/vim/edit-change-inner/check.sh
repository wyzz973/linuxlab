#!/bin/bash
expected='message = "goodbye world"
server = "production.server.com"
greeting = print(message)'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
