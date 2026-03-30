#!/bin/bash
expected='def process():
    if ready:
        result = compute()
        save(result)
    return True'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
