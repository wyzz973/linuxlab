#!/bin/bash
expected='.primary { color: #3498db; }
.danger { color: #e74c3c; }
.success { color: #2ecc71; }'
actual=$(cat challenge.txt 2>/dev/null)
[ "$actual" = "$expected" ]
