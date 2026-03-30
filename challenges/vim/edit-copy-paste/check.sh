#!/bin/bash
expected='.header {
  border: 1px solid #333;
  color: #fff;
  background: #000;
  border: 1px solid #333;
}'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
