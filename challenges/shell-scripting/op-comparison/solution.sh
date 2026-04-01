#!/bin/bash
a=15
b=20
[ "$a" -eq "$b" ] && echo "a == b: true" || echo "a == b: false"
[ "$a" -ne "$b" ] && echo "a != b: true" || echo "a != b: false"
[ "$a" -gt "$b" ] && echo "a > b: true" || echo "a > b: false"
[ "$a" -lt "$b" ] && echo "a < b: true" || echo "a < b: false"
[ "$a" -ge 10 ] && echo "a >= 10: true" || echo "a >= 10: false"
[ "$b" -le 20 ] && echo "b <= 20: true" || echo "b <= 20: false"
