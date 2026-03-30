#!/bin/bash
expected='[staging]
cache_enabled = no
debug_enabled = no
log_enabled = no
[production]
cache_enabled = yes
debug_enabled = yes
log_enabled = yes'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
