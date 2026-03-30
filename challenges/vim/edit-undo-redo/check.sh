#!/bin/bash
expected='level=info
level=warning
level=debug'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
