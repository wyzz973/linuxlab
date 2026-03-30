#!/bin/bash
expected='server:
  host: localhost
  port: 3000
database:
  name: mydb'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
