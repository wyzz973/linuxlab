#!/bin/bash
expected='def calculate_total(new_name):
    tax = new_name * 0.1
    discount = new_name * 0.05
    return new_name - discount + tax'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
